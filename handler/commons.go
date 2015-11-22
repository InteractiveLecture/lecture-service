package handler

import "net/http"

type errorHandler func(http.ResponseWriter, *http.Request) int

func (handler errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status := handler(w, r); status != -1 {
		w.WriteHeader(status)
	}
}

func createHandler(handlerFunc errorHandler) http.Handler {
	return http.Handler(errorHandler(handlerFunc))
}
