package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/*Валидация токена*/

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		c.Abort()
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		c.Abort()
		return
	}
	if len(headerParts[1]) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
		c.Abort()
		return
	}
	userId, err := h.services.ParseAccessToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error in parse": err.Error()})
		c.Abort()
		return
	}

	// Продолжаем выполнение запроса, т.к. Access Token действителен
	c.Set("userId", userId)
}
