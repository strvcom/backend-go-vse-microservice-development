package chi

import (
	"net/http"

	"vse-course/transport/model"
	"vse-course/transport/util"

	"github.com/go-chi/chi"
)

func getEmailFromURL(r *http.Request) string {
	email := chi.URLParam(r, "email")

	return email
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := util.UnmarshalRequest(r, &user)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = user.BirthDate.ValidateBirthDate()
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.Service.CreateUser(r.Context(), model.ToSvcUser(user))
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusCreated, user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.Service.GetUser(r.Context(), getEmailFromURL(r))
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, model.ToNetUser(user))
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Service.ListUsers(r.Context())

	util.WriteResponse(w, http.StatusOK, model.ToNetUsers(users))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := util.UnmarshalRequest(r, &user)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = user.BirthDate.ValidateBirthDate()
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	newUser, err := h.Service.UpdateUser(r.Context(), getEmailFromURL(r), model.ToSvcUser(user))
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, newUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.Service.DeleteUser(r.Context(), getEmailFromURL(r))
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusNoContent, nil)
}
