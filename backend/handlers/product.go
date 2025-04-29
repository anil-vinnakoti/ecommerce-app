package routes

import (
	"github.com/anil-vinnakoti/ecommerce-app/backend/db"
	"github.com/anil-vinnakoti/ecommerce-app/backend/models"
	"github.com/gin-gonic/gin"
)

// GetProducts returns all products
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := db.DB.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to fetch products"})
		return
	}
	c.JSON(200, products)
}

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	if err := db.DB.Create(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to create product"})
		return
	}
	c.JSON(201, product)
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	if err := db.DB.Save(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to update product"})
		return
	}
	c.JSON(200, product)
}

// DeleteProduct deletes an existing product
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to delete product"})
		return
	}
	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
