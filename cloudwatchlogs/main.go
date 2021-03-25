package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"encoding/json"

	"github.com/DavidGamba/go-getoptions"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

var healthy int32
var done chan bool
var quit chan os.Signal
var myVersion string = "0.0.1"

//StatusRecorder records the status sent to the client
type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

// WriteHeader is a wrapper so we can capture the status code.
func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// HealthCheckBody struct
type HealthCheckBody struct {
	Alive   int32
	Version string
	OK      bool `json:"ok"`
}

func index(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "bad") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `Valid Routes
/               index - help
/health   endpoint to check status of the service
/bad/...  forces a 404
-------------------------------
`)
	fmt.Fprintln(w, "Request Info:")
	fmt.Fprintln(w, "URL Path:", r.URL.Path)
	fmt.Fprintln(w, "Host:", r.Host)

	fmt.Fprintln(w, "\nHeaders Received")
	for key, val := range r.Header {
		fmt.Fprintf(w, "%s: %v\n", key, val)
	}

	fmt.Fprintln(w, "\nQuery Params")
	for key, val := range r.URL.Query() {
		fmt.Fprintf(w, "%s: %v\n", key, val)
	}
}

// HandleHealthCheck controller
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {

	myhealth := HealthCheckBody{healthy, myVersion, true}

	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(myhealth)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func customLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(recorder, r)

		myFields := log.Fields{
			"status":           recorder.Status,
			"method":           r.Method,
			"path":             r.URL.Path,
			"ip":               strings.Split(r.RemoteAddr, ":")[0],
			"user_agent":       r.UserAgent(),
			"content_encoding": r.Header.Get("Content-Encoding"),
		}

		myQueryFields := log.Fields{}
		for key, val := range r.URL.Query() {
			myQueryFields[key] = val
		}

		myFields["queryParams"] = myQueryFields
		myHeaderFields := log.Fields{}
		for key, val := range r.Header {
			myHeaderFields[key] = val
		}

		myFields["headers"] = myHeaderFields
		log.WithFields(myFields).Info()
	})
}

func main() {
	os.Exit(runProgram(os.Args[1:]))
}

func runProgram(args []string) int {

	done = make(chan bool)
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	opt := getoptions.New()
	opt.Bool("help", false, opt.Alias("h", "?"))
	opt.Bool("version", false, opt.Alias("V"))
	opt.Bool("disable-json", false, opt.GetEnv("DISABLE_JSON"), opt.Description("Disable JSON logging"))
	remaining, err := opt.Parse(args)
	if opt.Called("help") {
		fmt.Fprint(os.Stderr, opt.Help())
		return 0
	}
	if opt.Called("version") {
		fmt.Println("Version:", myVersion)
		return 0
	}
	if !opt.Called("disable-json") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		fmt.Fprint(os.Stderr, opt.Help(getoptions.HelpSynopsis))
		return 2
	}

	if len(remaining) > 0 {
		log.WithFields(
			log.Fields{
				"Err": remaining,
			},
		).Error("Unparsed options")
	}
	httpListenAddress := ":8080"

	log.Info("Starting app")

	router := mux.NewRouter()
	router.HandleFunc("/health", HandleHealthCheck)
	router.PathPrefix("/").HandlerFunc(index)
	router.Use(handlers.ProxyHeaders)
	router.Use(customLogger)

	//We are using a collection of handlers from here:
	//https://www.gorillatoolkit.org/pkg/handlers

	server := &http.Server{
		Addr:         httpListenAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  3 * time.Second,
	}

	go func() {
		<-quit
		signal.Stop(quit)
		log.Info("Server is shutting down")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.WithFields(
				log.Fields{
					"Err": err,
				},
			).Fatal("Could not gracefull shutdown the server")
		}
		done <- true
	}()

	log.Infof("Server is ready to handle requests at: %v", httpListenAddress)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithFields(
			log.Fields{
				"Err": err,
			},
		).Fatalf("Could not listen on %s", httpListenAddress)
	}

	<-done
	log.Info("Server stopped")
	return 0
}
