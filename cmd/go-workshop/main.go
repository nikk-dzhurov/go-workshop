package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Server error: %s\n", err.Error())
	}
}
