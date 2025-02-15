package handlers

import (
	"net/http"

	"avito-tech-go/internal/services"
	"github.com/gin-gonic/gin"
)

func BuyMerchHandler(merchService services.MerchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "unauthorized"})
			return
		}

		itemType := c.Param("item")
		if itemType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "item type not specified"})
			return
		}

		err := merchService.BuyItem(userID.(uint), itemType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Товар успешно куплен",
			"item":    itemType,
		})
	}
}
