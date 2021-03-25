package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIndexHandler(t *testing.T) {

	type args struct {
		res *httptest.ResponseRecorder
		req *http.Request
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req404, _ := http.NewRequest("GET", "/bad/page", nil)

	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{"Index Page", args{httptest.NewRecorder(), req}, 200},
		{"Missing Page", args{httptest.NewRecorder(), req404}, 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index(tt.args.res, tt.args.req)
			assert.Equal(t, tt.statusCode, tt.args.res.Code)
		})
	}

}

func TestHandleHealthCheck(t *testing.T) {
	type args struct {
		res *httptest.ResponseRecorder
		req *http.Request
	}

	atomic.StoreInt32(&healthy, 1)
	HealthGET, _ := http.NewRequest("GET", "/health", nil)

	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{"GET Action", args{httptest.NewRecorder(), HealthGET}, 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleHealthCheck(tt.args.res, tt.args.req)
			assert.Equal(t, tt.statusCode, tt.args.res.Code)
		})
	}
}

//Test_runProgramNoServer should not start the HTTP server as it will block.
func Test_runProgramNoServer(t *testing.T) {

	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Help", args{args: []string{"--help"}}, 0},
		{"Version", args{args: []string{"--version"}}, 0},
		{"WrongArgs", args{args: []string{"--sfdsfsdfsdf"}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := runProgram(tt.args.args); got != tt.want {
				t.Errorf("runProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runProgramStartServer(t *testing.T) {

	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"HTTPPort", args{args: []string{}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			go runProgram(tt.args.args)

			select {
			case status := <-done:
				fmt.Println("TESTS: Sever shutdown completed:", status)
			case <-time.After(2 * time.Second):
				myPid := os.Getpid()
				fmt.Println("TESTS: Sending SIGINT to PID:", myPid)
				_ = syscall.Kill(myPid, syscall.SIGINT)
			}

		})
	}
}
