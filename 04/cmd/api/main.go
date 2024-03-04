package main

import (
	"context"
	"fmt"
	"log/slog"

	"user-management-api/service"
	"user-management-api/transport/api"
	"user-management-api/transport/util"

	httpx "go.strv.io/net/http"
)

var version = "v0.0.0"

func main() {
	ctx := context.Background()
	cfg := MustLoadConfig()
	util.SetServerLogLevel(slog.LevelInfo)

	controller, err := setupController(
		cfg,
	)
	if err != nil {
		slog.Error("initializing controller", slog.Any("error", err))
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	// Initialize the server config.
	serverConfig := httpx.ServerConfig{
		Addr:    addr,
		Handler: controller,
		Hooks:   httpx.ServerHooks{
			// BeforeShutdown: []httpx.ServerHookFunc{
			// 	func(_ context.Context) {
			// 		database.Close()
			// 	},
			// },
		},
		Limits: nil,
		Logger: util.NewServerLogger("httpx.Server"),
	}
	server := httpx.NewServer(&serverConfig)

	slog.Info("starting server", slog.Int("port", cfg.Port))
	if err := server.Run(ctx); err != nil {
		slog.Error("server failed", slog.Any("error", err))
	}
}

func setupController(
	_ Config,
) (*api.Controller, error) {
	// Initialize the service.
	svc, err := service.NewService()
	if err != nil {
		return nil, fmt.Errorf("initializing user service: %w", err)
	}

	// Initialize the controller.
	controller, err := api.NewController(
		svc,
		version,
	)
	if err != nil {
		return nil, fmt.Errorf("initializing controller: %w", err)
	}

	return controller, nil
}
