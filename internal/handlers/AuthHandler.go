package handlers

import (
	"avito-tech-go/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthRequest represents the request payload for authentication.
// swagger:model AuthRequest
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response containing the JWT token.
// swagger:model AuthResponse
type AuthResponse struct {
	Token string `json:"token"`
}

// AuthHandler godoc
// @Summary      Authenticate user and return JWT token
// @Description  If the user does not exist, the service registers the user and returns a token; otherwise, it performs login.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      AuthRequest  true  "Authentication request payload"
// @Success      200   {object}  AuthResponse
// @Failure      400   {object}  map[string]string "Invalid request payload"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Router       /api/auth [post]
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
