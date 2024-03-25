package v1

import (
	"net/http"
	"user-management-api/pkg/id"
	"user-management-api/transport/util"

	"github.com/go-chi/chi"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - CreateUser")
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	var userID id.User
	if err := userID.FromString(chi.URLParam(r, "id")); err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
	}
	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	util.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}
	util.WriteResponse(w, http.StatusOK, users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - UpdateUser")
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented - DeleteUser")
}
