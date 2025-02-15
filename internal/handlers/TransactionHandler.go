package handlers

import (
	"net/http"

	"avito-tech-go/internal/repositories"
	"avito-tech-go/internal/services"
	"github.com/gin-gonic/gin"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

func SendCoinHandler(txService services.TransactionService, userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SendCoinRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "invalid JSON request"})
			return
		}

		fromUserID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "unauthorized"})
			return
		}

		toUser, err := userRepo.GetUserByName(req.ToUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
			return
		}
		if toUser == nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "target user not found"})
			return
		}

		err = txService.TransferCoins(fromUserID.(uint), toUser.ID, req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Монеты успешно отправлены",
			"toUser":  req.ToUser,
			"amount":  req.Amount,
		})
	}
}
