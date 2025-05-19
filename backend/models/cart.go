// models/cart.go
package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    uint       `json:"user_id"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}
