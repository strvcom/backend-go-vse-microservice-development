package v1

import "github.com/go-chi/chi"

type Handler struct {
	*chi.Mux

	service Service
}

func NewHandler(
	service Service,
) *Handler {
	h := &Handler{
		service: service,
	}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {
	r := chi.NewRouter()

	// TODO: Setup middleware.

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.ListUsers)
		r.Post("/", h.CreateUser)
		r.Get("/{email}", h.GetUser)
		r.Put("/{email}", h.UpdateUser)
		r.Delete("/{email}", h.DeleteUser)
	})
}
