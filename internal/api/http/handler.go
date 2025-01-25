package apiHttp

import (
	"log/slog"
	"net/http"

	sloggin "github.com/samber/slog-gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mestvl-shop-app/backend/docs"
	"github.com/mestvl-shop-app/backend/pkg/auth"
	"github.com/mestvl-shop-app/backend/pkg/limiter"
	"github.com/mestvl-shop-app/backend/pkg/validator"

	appV1 "github.com/mestvl-shop-app/backend/internal/api/http/app/v1"
	"github.com/mestvl-shop-app/backend/internal/config"
	"github.com/mestvl-shop-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	logger       *slog.Logger
	tokenManager auth.TokenManager
}

func NewHandlers(services *service.Services, logger *slog.Logger, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		logger:       logger,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	validator.RegisterGinValidator()

	router.Use(
		gin.Recovery(),
		sloggin.NewWithConfig(h.logger, sloggin.Config{
			WithSpanID:  true,
			WithTraceID: true,
		}),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL, h.logger),
		corsMiddleware,
	)

	if cfg.HttpServer.SwaggerEnabled {
		h.logger.Info("swagger enabled")

		router.GET("/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/swagger/app/index.html")
		})

		router.GET("/swagger/app", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/swagger/app/index.html")
		})

		router.GET("/swagger/app/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("appApiV1")))
	}

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	appHandlersV1 := appV1.NewHandler(h.services, h.logger, h.tokenManager)
	api := router.Group("/api")
	{
		appHandlersV1.Init(api)
	}
}
