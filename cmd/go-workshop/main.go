package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nikk-dzhurov/go_workshop/internal/diagnostics"
)

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

	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/", helloHandler)
		server := &http.Server{
			Addr:    ":" + blPort,
			Handler: router,
		}

		log.Printf("The application server is starting on port: %s\n", blPort)
		err := server.ListenAndServe()
		if err != nil {
			possibleErrors <- fmt.Errorf("Server error: %s\n", err.Error())
		}
	}()

	go func() {
		diagnosticsRouter := diagnostics.NewDiagnostics()
		diagServer := &http.Server{
			Addr:    ":" + diagnosticsPort,
			Handler: diagnosticsRouter,
		}

		log.Printf("The diagnostics server is starting on port: %s\n", diagnosticsPort)
		err := diagServer.ListenAndServe()
		if err != nil {
			possibleErrors <- fmt.Errorf("Diagnostics server error: %s\n", err.Error())
		}
	}()

	select {
	case err := <-possibleErrors:
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("The hello handler was called")
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
