package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/mestvl-shop-app/backend/internal/domain"
	"github.com/mestvl-shop-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initClientRoutes(api *gin.RouterGroup) {
	clients := api.Group("/clients")
	clients.POST("/register", h.clientRegister)
	clients.POST("/login", h.clientLogin)
	clients.POST("/ping", h.clientIdentityMiddleware, h.ping)
}

type clientRegisterRequest struct {
	Email     string                     `json:"email" binding:"required,email"`
	Password  string                     `json:"password" binding:"required,min=6"`
	Firstname string                     `json:"firstname" binding:"required"`
	Surname   string                     `json:"surname" binding:"required"`
	Birthday  *time.Time                 `json:"birthday" binding:"omitempty,datetime"`
	Gender    *domain.ClientGenderString `json:"gender" binding:"omitempty"`
}

// @Summary Регистрация
// @Tags Client
// @Description Регистрация
// @ModuleID Client
// @Accept  json
// @Produce  json
// @Param input body clientRegisterRequest true "Регистрация"
// @Success 201
// @Failure 400 {object} ErrorStruct
// @Failure 500
// @Router /clients/register [post]
func (h *Handler) clientRegister(c *gin.Context) {
	var req clientRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	err := h.services.Client.Register(c.Request.Context(), &service.RegisterClientInput{
		Email:     req.Email,
		Password:  req.Password,
		Firstname: req.Firstname,
		Surname:   req.Surname,
		Birthday:  req.Birthday,
		Gender:    req.Gender,
	})

	if err != nil {
		if errors.Is(err, service.ClientAlreadyExists) {
			errorResponse(c, ClientAlreadyExistsCode)
			return
		}

		h.logger.Error("failed to register client",
			"error", err,
		)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusCreated)
}

type clientLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type clientLoginResponse struct {
	AccessToken string `json:"access_token"`
}

// @Summary Авторизация
// @Tags Client
// @Description Авторизация
// @ModuleID Client
// @Accept  json
// @Produce  json
// @Param input body clientLoginRequest true "Авторизация"
// @Success 200 {object} clientLoginResponse
// @Failure 400 {object} ErrorStruct
// @Failure 500
// @Router /clients/login [post]
func (h *Handler) clientLogin(c *gin.Context) {
	var req clientLoginRequest
	if err := c.BindJSON(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	token, err := h.services.Client.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, service.ClientInvalidCredentials) {
			errorResponse(c, ClientNotFoundCode)
			return
		}

		h.logger.Error("failed to login client",
			"error", err,
		)
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, clientLoginResponse{AccessToken: token})
}

// @Summary Ping
// @Tags Client
// @Description Ping
// @ModuleID Client
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 403
// @Router /clients/ping [post]
// @Security UserAuth
func (h *Handler) ping(c *gin.Context) {
	c.Status(http.StatusOK)
}
