package utils

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SecureRoute(router *mux.Router, path string, handler http.HandlerFunc, method string) {
	router.HandleFunc(path, ValidateJWT(handler)).Methods(method)
}
