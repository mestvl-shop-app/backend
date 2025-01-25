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
	clients := api.Group("/clients/register")
	clients.POST("", h.clientRegister)
}

type clientRegisterRequest struct {
	Firstname string                     `json:"firstname" binding:"required"`
	Surname   string                     `json:"surname" binding:"required"`
	Birthday  *time.Time                 `json:"birthday" binding:"omitempty"`
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
// @Router /clients [post]
func (h *Handler) clientRegister(c *gin.Context) {
	var req clientRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		validationErrorResponse(c, err)
		return
	}

	err := h.services.Client.Register(c.Request.Context(), &service.RegisterClientDTO{
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
