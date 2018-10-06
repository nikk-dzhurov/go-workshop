package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("Server error: %s\n", err.Error())
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
