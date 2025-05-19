package routes

import (
	"net/http"

	"github.com/anil-vinnakoti/ecommerce-app/backend/db"
	"github.com/anil-vinnakoti/ecommerce-app/backend/models"
	"github.com/gin-gonic/gin"
)

func CreateCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func GetCart(c *gin.Context) {
	id := c.Param("id")
	var cart models.Cart
	if err := db.DB.Preload("CartItems.Product").First(&cart, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// Add item to cart
func AddToCart(c *gin.Context) {
	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	// get user from session
	userID := getSessionUserID(c)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItem models.CartItem
	db.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&cartItem)

	if cartItem.ID == 0 {
		cartItem = models.CartItem{UserID: userID, ProductID: input.ProductID, Quantity: input.Quantity}
		db.DB.Create(&cartItem)
	} else {
		cartItem.Quantity += input.Quantity
		db.DB.Save(&cartItem)
	}

	c.JSON(http.StatusOK, cartItem)
}

// Remove item from cart
func RemoveFromCart(c *gin.Context) {
	cartItemID := c.Param("id")
	db.DB.Delete(&models.CartItem{}, cartItemID)
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func CartTotal(c *gin.Context) {
	userID := 1

	var items []models.CartItem
	db.DB.Preload("Product").Where("user_id = ?", userID).Find(&items)

	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * item.Product.Price
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
