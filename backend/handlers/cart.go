package routes

import (
	"net/http"

	"github.com/anil-vinnakoti/ecommerce-app/backend/db"
	"github.com/anil-vinnakoti/ecommerce-app/backend/middleware"
	"github.com/anil-vinnakoti/ecommerce-app/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Helper to get user ID from context (set by JWT middleware)
func getUserID(c *gin.Context) uint {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return 0
	}
	return userIDVal.(uint)
}

func CreateCart(c *gin.Context) {
	userID := getUserID(c)

	cart := models.Cart{UserID: userID}
	if err := db.DB.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func GetCart(c *gin.Context) {
	userID := getUserID(c)

	var cart models.Cart
	err := db.DB.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// Add item to cart
func AddToCart(c *gin.Context) {
	// Assume you have a function to get authenticated user ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find or create a Cart for the user
	var cart models.Cart
	err := db.DB.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cart = models.Cart{UserID: userID}
			if err := db.DB.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// Check if CartItem for this product already exists in this cart
	var cartItem models.CartItem
	err = db.DB.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&cartItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new CartItem
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: input.ProductID,
				Quantity:  input.Quantity,
			}
			if err := db.DB.Create(&cartItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		// Update quantity if item exists
		cartItem.Quantity += input.Quantity
		if err := db.DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
	}

	c.JSON(http.StatusOK, cartItem)
}

// Remove item from cart
func RemoveFromCart(c *gin.Context) {
	cartItemID := c.Param("id")
	userID := getUserID(c)

	var cartItem models.CartItem
	if err := db.DB.First(&cartItem, "id = ? AND user_id = ?", cartItemID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found or access denied"})
		return
	}

	if err := db.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func CartTotal(c *gin.Context) {
	userID := getUserID(c)

	var cart models.Cart
	db.DB.Where("user_id = ?", userID).First(&cart)

	var items []models.CartItem
	db.DB.Preload("Product").Where("cart_id = ?", cart.ID).Find(&items)

	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * item.Product.Price
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
