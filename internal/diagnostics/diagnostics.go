package diagnostics

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewDiagnostics() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ready", readyHandler)
	router.HandleFunc("/healthz", healthzHandler)

	return router
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, http.StatusText(http.StatusOK))
}
