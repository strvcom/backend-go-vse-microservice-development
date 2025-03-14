package v1

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	*chi.Mux

	service   Service
	validator *validator.Validate
}

func NewHandler(
	service Service,
) *Handler {
	h := &Handler{
		Mux:       chi.NewRouter(),
		service:   service,
		validator: validator.New(),
	}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {

	h.Route("/users", func(r chi.Router) {
		r.Get("/", h.ListUsers)
		r.Post("/", h.CreateUser)
		r.Get("/{email}", h.GetUser)
		r.Put("/{email}", h.UpdateUser)
		r.Delete("/{email}", h.DeleteUser)
	})
}
