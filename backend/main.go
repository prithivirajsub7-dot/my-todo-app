// main.go
package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// ✅ Database connect
	// ConnectDB() should return (*gorm.DB, error) and assign to global DB variable
	if err := ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// ✅ Initialize Gin router
	r := gin.Default()

	// ✅ CORS Middleware (Frontend ports 3000 & 3001)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Test route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server Running"})
	})

	// 🔐 Auth routes
	r.POST("/register", Register)
	r.POST("/login", Login)

	// 🔒 Protected routes
	protected := r.Group("/")
	protected.Use(AuthMiddleware()) // Auth middleware sets userID in context
	{
		protected.POST("/todos", CreateTodo)
		protected.GET("/todos", GetTodos)
		protected.PUT("/todos/:id", UpdateTodo)
		protected.DELETE("/todos/:id", DeleteTodo)
	}

	// ✅ Start server with logging
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}