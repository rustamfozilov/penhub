package main

import (
	"github.com/rustamfozilov/penhub/internal/handlers"
	"net/http"
)

func NewRouter(h *handlers.Handler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/createbook", h.CreateBook)


	return mux
}
