package nethttp

import (
	"net/http"

	"vse-course/service"
	"vse-course/transport/model"
)

type Handler struct {
	Port    int
	Mux     *http.ServeMux
	Service model.Service
}

func Initialize(port int) *Handler {
	h := &Handler{
		Port:    port,
		Mux:     http.NewServeMux(),
		Service: service.CreateService(),
	}

	h.Mux.HandleFunc("/users", h.HandleUsers)
	h.Mux.HandleFunc("/users/", h.HandleUser)

	return h
}
