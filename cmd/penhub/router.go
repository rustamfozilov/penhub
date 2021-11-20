package main

import (
	"github.com/rustamfozilov/penhub/internal/handlers"
	"net/http"
)

func NewRouter(h *handlers.Handler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/createbook", h.CreateBook)
	mux.HandleFunc("/registration_user", h.RegistrationUser)
	mux.HandleFunc("/get_token_user", h.GetTokenForUser)

	return mux
}
