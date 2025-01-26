package v1

import (
	"log/slog"

	"github.com/mestvl-shop-app/backend/internal/client/auth"
	"github.com/mestvl-shop-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// @title Backend API
// @version 1.0
// @description Backend API for NNBlog Service
// @BasePath /api/app/v1
// @securityDefinitions.apikey UserAuth
// @in header
// @name Authorization

type Handler struct {
	services          *service.Services
	logger            *slog.Logger
	authServiceClient *auth.Client
}

func NewHandler(
	services *service.Services,
	logger *slog.Logger,
	authServiceClient *auth.Client,
) *Handler {
	return &Handler{
		services:          services,
		logger:            logger,
		authServiceClient: authServiceClient,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("app/v1")
	{
		h.initClientRoutes(v1)
	}
}
