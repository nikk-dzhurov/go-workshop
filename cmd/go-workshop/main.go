package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikk-dzhurov/go_workshop/internal/diagnostics"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler)
	// router.HandleFunc("/healthz", helloHandler)

	go func() {
		log.Println("Start server on port: 8080")
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatalln("Server error: %s\n", err.Error())
		}
	}()

	diagnostics := diagnostics.NewDiagnostics()

	log.Println("Start diagnostics server on port: 8585")
	err := http.ListenAndServe(":8585", diagnostics)
	if err != nil {
		log.Fatalln("Server error: %s\n", err.Error())
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
