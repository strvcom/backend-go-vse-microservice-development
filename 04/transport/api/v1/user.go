package v1

import (
	"net/http"

	"github.com/go-chi/chi"
)

func getEmailFromURL(r *http.Request) string {
	email := chi.URLParam(r, "email")
	return email
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - CreateUser")
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - GetUser")
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - ListUsers")
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - UpdateUser")
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - DeleteUser")
}
