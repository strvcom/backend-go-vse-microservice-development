package nethttp

import (
	"net/http"
	"regexp"

	svcmodel "vse-course/service/model"
	"vse-course/transport/model"
	"vse-course/transport/util"
)

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		util.WriteErrResponse(w, http.StatusMethodNotAllowed, nil)
	}
}

var usersRegExp = regexp.MustCompile(`/users/(?P<email>[\w-\.]+@([\w-]+\.)+[\w-]{2,4})`)

func (h *Handler) HandleUser(w http.ResponseWriter, r *http.Request) {
	match := usersRegExp.FindStringSubmatch(r.URL.Path)

	var email string
	if len(match) > 0 {
		for i, name := range usersRegExp.SubexpNames() {
			if name == "email" {
				email = match[i]
			}
		}
	} else {
		util.WriteErrResponse(w, http.StatusNotFound, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r, email)
		return
	case http.MethodPatch:
		h.UpdateUser(w, r, email)
		return
	case http.MethodDelete:
		h.DeleteUser(w, r, email)
	default:
		util.WriteErrResponse(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Service.ListUsers(r.Context())

	util.WriteResponse(w, http.StatusOK, model.ToNetUsers(users))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := util.UnmarshalRequest(r, &user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.Service.CreateUser(r.Context(), model.ToSvcUser(user)); err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, http.NoBody)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, email string) {
	user, err := h.Service.GetUser(r.Context(), email)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, email string) {
	var user model.User
	err := util.UnmarshalRequest(r, &user)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	newUser, err := h.Service.UpdateUser(r.Context(), email, svcmodel.User{})
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusOK, newUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, email string) {
	err := h.Service.DeleteUser(r.Context(), email)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusNoContent, nil)
}
