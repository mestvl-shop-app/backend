package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiHttp "github.com/mestvl-shop-app/backend/internal/api/http"
	"github.com/mestvl-shop-app/backend/internal/client/auth"
	"github.com/mestvl-shop-app/backend/internal/config"
	"github.com/mestvl-shop-app/backend/internal/db"
	"github.com/mestvl-shop-app/backend/internal/repository"
	"github.com/mestvl-shop-app/backend/internal/server"
	"github.com/mestvl-shop-app/backend/internal/service"

	log "github.com/mestvl-shop-app/backend/pkg/logger"
)

const (
	AppID = 1
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	logger := log.SetupLogger(cfg.Env)
	logger.Info("start shop backend",
		"env", cfg.Env,
	)
	logger.Debug("debug messages are enabled")

	// Init database
	dbPostgres, err := db.New(cfg.Database)
	if err != nil {
		logger.Error("postgres connect problem", "error", err)
		os.Exit(1)
	}
	defer func() {
		err = dbPostgres.Close()
		if err != nil {
			logger.Error("error when closing", "error", err)
		}
	}()
	logger.Info("postgres connection done")

	// Init auth service grpc client
	authServiceClient, err := auth.New(
		context.Background(),
		logger,
		cfg.Clients.AuthService.Address,
		cfg.Clients.AuthService.Timeout,
		cfg.Clients.AuthService.RetriesCount,
	)
	if err != nil {
		logger.Error("failed to init auth service client",
			"error", err,
		)
		return
	}

	// Init services, repositories, handlers
	repos := repository.NewRepositories(dbPostgres)
	services := service.NewServices(service.Deps{
		Logger:     logger,
		Config:     cfg,
		Repos:      repos,
		AuthClient: authServiceClient,
	})
	handlers := apiHttp.NewHandlers(
		services,
		logger,
		authServiceClient,
	)

	// Init HTTP server
	srv := server.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error occurred while running http server", "error", err)
		}
	}()
	logger.Info("server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Error("failed to stop server", "error", err)
	}

	logger.Info("app stopped")
}
