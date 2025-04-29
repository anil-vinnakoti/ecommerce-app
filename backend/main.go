package main

import (
	"log"

	"github.com/anil-vinnakoti/ecommerce-app/backend/db"
	routes "github.com/anil-vinnakoti/ecommerce-app/backend/handlers"
	"github.com/anil-vinnakoti/ecommerce-app/backend/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	db.Connect()

	// Auto-migrate the Product model to create the table
	err := db.DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}

	log.Println("Database migration complete!")

	// Set up Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Define routes for CRUD operations on Products
	r.GET("/products", routes.GetProducts)
	r.POST("/products", routes.CreateProduct)
	r.PUT("/products/:id", routes.UpdateProduct)
	r.DELETE("/products/:id", routes.DeleteProduct)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
