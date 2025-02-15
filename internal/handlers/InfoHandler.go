package handlers

import (
	"net/http"

	"avito-tech-go/internal/services"
	"github.com/gin-gonic/gin"
)

// InfoHandler godoc
// @Summary      Get user's coin info, inventory, and transaction history
// @Description  Retrieves the coin balance, purchased merch items, and coin transaction history for the authenticated user.
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  services.InfoResponse
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/info [get]
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
