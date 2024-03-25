package main

import (
	"context"
	"fmt"
	"log/slog"

	"user-management-api/repository"
	"user-management-api/service"
	"user-management-api/transport/api"
	"user-management-api/transport/util"

	"github.com/jackc/pgx/v5/pgxpool"
	httpx "go.strv.io/net/http"
)

var version = "v0.0.0"

func main() {
	ctx := context.Background()
	cfg := MustLoadConfig()
	util.SetServerLogLevel(slog.LevelInfo)

	database, err := setupDatabase(ctx, cfg)
	if err != nil {
		slog.Error("initializing database", slog.Any("error", err))
	}
	repository, err := repository.New(database)
	if err != nil {
		slog.Error("initializing repository", slog.Any("error", err))
	}

	controller, err := setupController(
		cfg,
		repository,
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

func setupDatabase(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	// Initialize the database connection pool.
	pool, err := pgxpool.New(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func setupController(
	_ Config,
	repository service.Repository,
) (*api.Controller, error) {
	// Initialize the service.
	svc, err := service.NewService(repository)
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
