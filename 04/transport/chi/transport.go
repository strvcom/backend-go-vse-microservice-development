package chi

import (
	"net/http"

	"vse-course/service"
	"vse-course/transport/model"

	"github.com/go-chi/chi"
)

type Handler struct {
	Port    int
	Mux     *chi.Mux
	Service model.Service
}

func Initialize(port int) *Handler {
	h := &Handler{
		Port:    port,
		Mux:     chi.NewRouter(),
		Service: service.CreateService(),
	}

	h.Mux.Get("/healthz", healthz)

	h.Mux.Route("/users", func(r chi.Router) {
		r.Get("/", h.ListUsers)
		r.Post("/", h.CreateUser)

		r.Route("/{email}", func(r chi.Router) {
			r.Get("/", h.GetUser)
			r.Delete("/", h.DeleteUser)
			r.Patch("/", h.UpdateUser)
		})
	})

	return h
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
