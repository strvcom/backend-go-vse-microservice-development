package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	"user-management-api/service/model"
	transportmodel "user-management-api/transport/api/v1/model"
	"user-management-api/transport/util"
)

func getEmailFromURL(r *http.Request) string {
	email := chi.URLParam(r, "email")
	return email
}

func (h *Handler) validateUser(user transportmodel.User) error {
	if err := h.validator.Struct(user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, e := range validationErrors {
				switch e.Tag() {
				case "required":
					errMsgs = append(errMsgs, fmt.Sprintf("%s is required", e.Field()))
				case "email":
					errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid email address", e.Field()))
				default:
					errMsgs = append(errMsgs, fmt.Sprintf("%s failed %s validation", e.Field(), e.Tag()))
				}
			}
			return fmt.Errorf("validation failed: %s", errMsgs)
		}
		return err
	}
	return nil
}

func mapTransportUserToServiceUser(user transportmodel.User) model.User {
	return model.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func mapServiceUserToTransportUser(user model.User) transportmodel.User {
	return transportmodel.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user transportmodel.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.validateUser(user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	serviceUser := mapTransportUserToServiceUser(user)
	if err := h.service.CreateUser(r.Context(), serviceUser); err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteResponse(w, http.StatusCreated, user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := getEmailFromURL(r)
	if email == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}
	
	user, err := h.service.GetUser(r.Context(), email)
	if err != nil {
		util.WriteErrResponse(w, http.StatusNotFound, err)
		return
	}

	transportUser := mapServiceUserToTransportUser(user)
	util.WriteResponse(w, http.StatusOK, transportUser)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.ListUsers(r.Context())
	
	// Pre-allocate with exact size for better performance
	transportUsers := make([]transportmodel.User, len(users))
	for i, user := range users {
		transportUsers[i] = mapServiceUserToTransportUser(user)
	}
	
	util.WriteResponse(w, http.StatusOK, transportUsers)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	email := getEmailFromURL(r)
	if email == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}

	var user transportmodel.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := h.validateUser(user); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	serviceUser := mapTransportUserToServiceUser(user)
	updatedUser, err := h.service.UpdateUser(r.Context(), email, serviceUser)
	if err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	transportUser := mapServiceUserToTransportUser(updatedUser)
	util.WriteResponse(w, http.StatusOK, transportUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := getEmailFromURL(r)
	if email == "" {
		util.WriteErrResponse(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}

	if err := h.service.DeleteUser(r.Context(), email); err != nil {
		util.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
