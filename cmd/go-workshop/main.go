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

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler)

	go func() {
		log.Printf("The application server is starting on port: %s\n", blPort)
		err := http.ListenAndServe(":"+blPort, router)
		if err != nil {
			log.Fatalln("Server error: %s\n", err.Error())
		}
	}()

	diagnostics := diagnostics.NewDiagnostics()

	log.Printf("The diagnostics server is starting on port: %s\n", diagnosticsPort)
	err := http.ListenAndServe(":"+diagnosticsPort, diagnostics)
	if err != nil {
		log.Fatalln("Server error: %s\n", err.Error())
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("The hello handler was called")
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
