package api

import (
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"

	apiv1 "user-management-api/transport/api/v1"
	"user-management-api/transport/middleware"

	"github.com/go-chi/chi"
	httpx "go.strv.io/net/http"
)

//go:embed openapi.yaml
var OpenAPI []byte

// Controller handles all /api endpoints.
// It is responsible for routing requests to appropriate handlers.
// Versioned endpoints are handled by subcontrollers.
type Controller struct {
	*chi.Mux

	authenticator middleware.Authenticator
	service       apiv1.Service
	version       string
}

func NewController(
	authenticator middleware.Authenticator,
	service apiv1.Service,
	version string,
) (*Controller, error) {
	controller := &Controller{
		authenticator: authenticator,
		service:       service,
		version:       version,
	}
	controller.initRouter()
	return controller, nil
}

func (c *Controller) initRouter() {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		v1Handler := apiv1.NewHandler(
			c.authenticator,
			c.service,
		)

		r.Route("/api", func(r chi.Router) {
			r.Get("/openapi.yaml", c.OpenAPI)
			r.Mount("/v1", v1Handler)
		})
	})

	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/version", c.Version)

	c.Mux = r
}

// TODO: Improve this handler.
func (c *Controller) OpenAPI(w http.ResponseWriter, _ *http.Request) {
	encodeFunc := func(w http.ResponseWriter, data any) error {
		d, ok := data.([]byte)
		if !ok {
			return fmt.Errorf("expected byte slice: got %T", data)
		}
		if _, err := w.Write(d); err != nil {
			return fmt.Errorf("writing openapi content: %w", err)
		}
		return nil
	}
	if err := httpx.WriteResponse(
		w,
		OpenAPI,
		http.StatusOK,
		httpx.WithEncodeFunc(encodeFunc),
		httpx.WithContentType(httpx.ApplicationYAML),
	); err != nil {
		slog.Error("writing response", slog.Any("error", err))
	}
}

// TODO: Improve this handler.
func (c *Controller) Version(w http.ResponseWriter, _ *http.Request) {
	if err := httpx.WriteResponse(
		w,
		c.version,
		http.StatusOK,
	); err != nil {
		slog.Error("writing response", slog.Any("error", err))
	}
}
