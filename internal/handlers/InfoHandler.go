package handlers

import (
	"net/http"

	"avito-tech-go/internal/services"
	"github.com/gin-gonic/gin"
)

func InfoHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "unauthorized"})
			return
		}

		info, err := userService.GetInfo(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
			return
		}

		c.JSON(http.StatusOK, info)
	}
}
