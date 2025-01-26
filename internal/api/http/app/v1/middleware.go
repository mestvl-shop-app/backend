package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	authorizationHeader = "Authorization"
	clientCtx           = "clientId"
)

func (h *Handler) clientIdentityMiddleware(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		h.logger.Error("invalid auth header")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if len(headerParts[1]) == 0 {
		h.logger.Error("token is empty")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ok, err := h.authServiceClient.Validate(c.Request.Context(), headerParts[1])
	if err != nil {
		h.logger.Error("validate token failed",
			"error", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, _ := jwt.Parse(headerParts[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.logger.Error("get token claims failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var id string
	if claims["uid"] != nil {
		id = claims["uid"].(string)
	}

	c.Set(clientCtx, id)
}
