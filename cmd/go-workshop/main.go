package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nikk-dzhurov/go_workshop/internal/diagnostics"
)

type ServerConfig struct {
	port   string
	router http.Handler
	name   string
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

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler)
	diagnosticsRouter := diagnostics.NewDiagnostics()

	configs := []ServerConfig{
		{port: blPort, router: router, name: "Application server"},
		{port: diagnosticsPort, router: diagnosticsRouter, name: "Diagnostics server"},
	}

	servers := make([]*http.Server, 2)
	for i, c := range configs {
		go func(config ServerConfig, idx int) {
			servers[idx] = &http.Server{
				Addr:    ":" + config.port,
				Handler: config.router,
			}

			log.Printf("The %s is starting on port: %s\n", config.name, config.port)
			err := servers[idx].ListenAndServe()
			if err != nil {
				possibleErrors <- fmt.Errorf("%s error: %s\n", config.name, err.Error())
			}
		}(c, i)
	}

	select {
	case err := <-possibleErrors:
		for _, s := range servers {
			s.Shutdown(context.Background())
		}
		log.Fatal(err.Error())
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("The hello handler was called")
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
