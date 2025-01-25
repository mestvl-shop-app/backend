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
	"github.com/mestvl-shop-app/backend/internal/config"
	"github.com/mestvl-shop-app/backend/internal/db"
	"github.com/mestvl-shop-app/backend/internal/repository"
	"github.com/mestvl-shop-app/backend/internal/server"
	"github.com/mestvl-shop-app/backend/internal/service"
	"github.com/mestvl-shop-app/backend/pkg/auth"

	hash "github.com/mestvl-shop-app/backend/pkg/hasher"
	log "github.com/mestvl-shop-app/backend/pkg/logger"
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

	// Token manager
	tokenManager, err := auth.NewManager(cfg.Auth.JWT)
	if err != nil {
		logger.Error("auth manager creation err", "error", err)
		return
	}

	// Hasher
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	// Init services, repositories, handlers
	repos := repository.NewRepositories(dbPostgres)
	services := service.NewServices(service.Deps{
		Logger:       logger,
		Config:       cfg,
		Hasher:       hasher,
		TokenManager: tokenManager,
		Repos:        repos,
	})
	handlers := apiHttp.NewHandlers(services, logger, tokenManager)

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
