package handlers

import (
	"net/http"

	"avito-tech-go/internal/services"
	"github.com/gin-gonic/gin"
)

// BuyMerchHandler godoc
// @Summary      Purchase a merchandise item using coins
// @Description  Allows the authenticated user to buy a merch item specified by the item type.
// @Tags         merch
// @Security     BearerAuth
// @Produce      json
// @Param        item  path      string  true  "Merch item type"
// @Success      200   {object}  map[string]interface{} "Successful purchase response"
// @Failure      400   {object}  map[string]string "Bad request"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Router       /api/buy/{item} [get]
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
