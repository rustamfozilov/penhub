package handlers

import (
	"encoding/json"
	"errors"
	"github.com/rustamfozilov/penhub/internal/services"
	"github.com/rustamfozilov/penhub/internal/types"
	"net/http"
)

func (h *Handler) GetTokenForUser(w http.ResponseWriter, r *http.Request) {
	var u types.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		badRequest(w, err)
		return
	}

	token, err := h.Service.GetTokenForUser(r.Context(), &u)
	if errors.Is(err, services.ErrNoSuchUser) || errors.Is(err, services.ErrInvalidPassword) {
		badRequest(w, err)
		return
	}
	if err != nil {
		InternalServerError(w, err)
		return
	}

	item := types.T{Token: token}
	data, err := json.Marshal(item)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		InternalServerError(w, err)
		return
	}
}
