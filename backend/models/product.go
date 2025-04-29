package models

import (
	"gorm.io/gorm"
)

// Product represents an e-commerce product
type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"` // You can include an image URL for the product
}
