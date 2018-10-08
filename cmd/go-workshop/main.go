package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/nikk-dzhurov/go-workshop/internal/diagnostics"
)

type serverConfig struct {
	port   string
	router http.Handler
	name   string
}

type apiController struct {
	counter int
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	blPort := os.Getenv("PORT")
	if blPort == "" {
		log.Fatal("Provide PORT environment variable")
	}
	diagnosticsPort := os.Getenv("DIAG_PORT")
	if diagnosticsPort == "" {
		log.Fatal("Provide DIAG_PORT environment variable")
	}

	possibleErrors := make(chan error, 2)

	ctrl := &apiController{counter: 0}

	router := mux.NewRouter()
	router.HandleFunc("/", ctrl.helloHandler)
	diagnosticsRouter := diagnostics.NewDiagnostics()

	configs := []serverConfig{
		{port: blPort, router: router, name: "Application server"},
		{port: diagnosticsPort, router: diagnosticsRouter, name: "Diagnostics server"},
	}

	servers := make([]*http.Server, len(configs))
	for i, c := range configs {
		go func(config serverConfig, idx int) {
			servers[idx] = &http.Server{
				Addr:    ":" + config.port,
				Handler: config.router,
			}

			log.Printf("The %s is starting on port: %s\n", config.name, config.port)
			err := servers[idx].ListenAndServe()
			if err != nil {
				possibleErrors <- fmt.Errorf("%s error: %s", config.name, err.Error())
			}
		}(c, i)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-possibleErrors:
		log.Printf("Got an error: %s\n", err.Error())
	case sig := <-interrupt:
		log.Printf("Received the signal %v\n", sig)
	}

	for _, s := range servers {
		timeout := 5 * time.Second
		log.Printf("Shutdown with timeout: %s\n", timeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		customErr := s.Shutdown(ctx)
		if customErr != nil {
			log.Println(customErr.Error())
		}

		log.Printf("Server on address %s, gracefully stopped\n", s.Addr)
	}
}

func (ctrl *apiController) helloHandler(w http.ResponseWriter, r *http.Request) {
	ctrl.counter++
	log.Printf("The hello handler was called, count %d\n", ctrl.counter)
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
