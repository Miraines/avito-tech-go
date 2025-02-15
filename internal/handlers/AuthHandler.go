package handlers

import (
	"avito-tech-go/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func AuthHandler(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "invalid request payload"})
			return
		}

		loginToken, err := authService.Login(req.Username, req.Password)
		if err != nil {
			regToken, regErr := authService.Register(req.Username, req.Password)
			if regErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"errors": regErr.Error()})
				return
			}
			c.JSON(http.StatusOK, AuthResponse{Token: regToken})
			return
		}
		c.JSON(http.StatusOK, AuthResponse{Token: loginToken})
	}
}
