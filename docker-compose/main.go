package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

var dockerComposeBody string = `
services:
  testservice:
    image: ubuntu:latest
    environment:
      - COMPOSE_TEST=TestMe123
    command: sleep infinity
    volumes:
      - /tmp:/tmp
    init: true
  testservice2:
    build: .
`

func main() {

	fmt.Println("Example docker-compose management")
	fmt.Println("==================================")
	fmt.Println("")
	fmt.Println("Docker-compose body:")
	fmt.Println(dockerComposeBody)
	fmt.Println("")
	fmt.Println("==================================")

	ctx := context.TODO()

	p := createDockerProject(ctx, dockerComposeBody)

	fmt.Println("Docker compose project created, list of known services:")
	for _, service := range p.Services {
		fmt.Println("Service:", service.Name)
		fmt.Println("  Image:", service.Image)
	}

	srv, err := createDockerService()

	if err != nil {
		log.Fatalln("error create docker service:", err)
	}

	fmt.Println("Docker service up...")

	var buildOptions *api.BuildOptions
	buildOptions = &api.BuildOptions{}
	buildOptions.Services = p.ServiceNames()
	buildOptions.Deps = true

	createOpts := api.CreateOptions{
		Build:         buildOptions,
		Services:      p.ServiceNames(),
		RemoveOrphans: true,
		IgnoreOrphans: true,
	}

	err = srv.Up(ctx, p, api.UpOptions{Create: createOpts})
	if err != nil {
		log.Fatalln("error up:", err)
	}

	myExec(ctx, srv, p, "echo hello world")
	myExec(ctx, srv, p, "cd / && ls")

	fmt.Println("Docker service down...")
	err = srv.Down(ctx, p.Name, api.DownOptions{})
	if err != nil {
		log.Fatalln("error down:", err)
	}

}

func createDockerProject(ctx context.Context, data string) *types.Project {

	defaultDockerContext, _ := os.Getwd()

	configDetails := types.ConfigDetails{
		// Fake path, doesn't need to exist.
		WorkingDir: defaultDockerContext,
		ConfigFiles: []types.ConfigFile{
			{Filename: "docker-compose.yaml", Content: []byte(dockerComposeBody)},
		},
		Environment: nil,
	}

	projectName := "testproject"

	p, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		options.SetProjectName(projectName, true)
	})

	if err != nil {
		log.Fatalln("error load:", err)
	}
	addServiceLabels(p)
	return p
}

// createDockerService creates a docker service which can be
// used to interact with docker-compose.
func createDockerService() (api.Compose, error) {
	var srv api.Compose
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return srv, err
	}

	dockerContext := "default"

	//Magic line to fix error:
	//Failed to initialize: unable to resolve docker endpoint: no context store initialized
	myOpts := &flags.ClientOptions{Context: dockerContext, LogLevel: "error"}
	err = dockerCli.Initialize(myOpts)
	if err != nil {
		return srv, err
	}

	srv = compose.NewComposeService(dockerCli)

	return srv, nil
}

func myExec(ctx context.Context, srv api.Compose, p *types.Project, cmd string) {
	result, err := srv.Exec(ctx, p.Name, api.RunOptions{
		Service:     "testservice",
		Command:     []string{"/bin/bash", "-c", cmd},
		WorkingDir:  "/bin",
		Tty:         true,
		Environment: []string{"TONE=test1"},
	})
	log.Println("Command result:", result, " and err:", err)
}

/*
addServiceLabels adds the labels docker compose expects to exist on services.
This is required for future compose operations to work, such as finding
containers that are part of a service.
*/
func addServiceLabels(project *types.Project) {
	for i, s := range project.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  "/",
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False", // default, will be overridden by `run` command
		}
		project.Services[i] = s
	}
}
