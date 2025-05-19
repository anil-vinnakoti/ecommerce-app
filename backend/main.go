package main

import (
	"log"

	"github.com/anil-vinnakoti/ecommerce-app/backend/db"
	routes "github.com/anil-vinnakoti/ecommerce-app/backend/handlers"
	"github.com/anil-vinnakoti/ecommerce-app/backend/middleware"
	"github.com/anil-vinnakoti/ecommerce-app/backend/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	db.Connect()

	// Auto-migrate the models
	err := db.DB.AutoMigrate(&models.Product{}, &models.Cart{}, &models.CartItem{})
	if err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}

	log.Println("Database migration complete!")

	// Set up Gin router
	r := gin.Default()

	// Configure CORS for frontend running on localhost:3000
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Public product routes
	r.GET("/products", routes.GetProducts)
	r.GET("/products/:productId", routes.GetProduct)
	r.POST("/products", routes.CreateProduct)
	r.PUT("/products/:productId", routes.UpdateProduct)
	r.DELETE("/products/:productId", routes.DeleteProduct)

	// User login and signup
	r.POST("/signup", routes.Signup)
	r.POST("/login", routes.Login)

	// Group for routes requiring authentication
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())

	// Protected cart routes
	auth.POST("/cart", routes.AddToCart)
	auth.GET("/cart", routes.GetCart)
	auth.DELETE("/cart/item/:id", routes.RemoveFromCart)
	auth.GET("/cart/total", routes.CartTotal)

	// Run the server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
