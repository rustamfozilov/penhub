package handlers

import (
	"encoding/json"
	"errors"
	"github.com/rustamfozilov/penhub/internal/services"
	"github.com/rustamfozilov/penhub/internal/types"
	"net/http"
)

func (h *Handler) RegistrationUser(w http.ResponseWriter, r *http.Request) {
	var u types.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = h.Service.RegistrationUser(r.Context(), &u)
	if errors.Is(err, services.ErrLoginUsed) {
		badRequest(w, err)
		return
	}
	if err != nil {
		InternalServerError(w, err)
		return
	}

}
