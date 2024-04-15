package v1

import (
	"user-management-api/transport/middleware"

	"github.com/go-chi/chi"
)

type Handler struct {
	*chi.Mux

	authenticator middleware.Authenticator
	service       Service
}

func NewHandler(
	authenticator middleware.Authenticator,
	service Service,
) *Handler {
	h := &Handler{
		authenticator: authenticator,
		service:       service,
	}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {
	r := chi.NewRouter()

	// TODO: Setup middleware.
	authenticate := middleware.NewAuthenticate(h.authenticator)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.With(authenticate).Get("/", h.ListUsers)
		r.With(authenticate).Get("/{email}", h.GetUser)
		r.With(authenticate).Put("/{email}", h.UpdateUser)
		r.With(authenticate).Delete("/{email}", h.DeleteUser)
	})
	h.Mux = r
}
