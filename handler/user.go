package handler

import (
	"net/http"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/service"
)

type UserHandler struct {
	usr service.AuthService
}

func NewUserHandler(usr service.AuthService) *UserHandler {
	return &UserHandler{usr: usr}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	if !decodeJSON(w, r, &u) {
		return
	}

	if err := h.usr.Register(u.Name, u.Email, u.PasswordHash); err != nil {
		respondError(w, err)
		return
	}
	response := map[string]any{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
	}
	writeJSON(w, http.StatusCreated, response)
}


