package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/delta/middleware"
	"github.com/go-Echelon/go-Echelon/pkg/delta/routes"
)

func main() {
	// TODO: Load from Echelon.toml config
	mongoURI := "mongodb://localhost:27017"
	dbName := "echelon_db"
	port := "8080"

	log.Println("🔌 Connecting to MongoDB...")

	db, err := database.Connect(mongoURI, dbName)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	log.Printf("✅ MongoDB connected: %s/%s\n", mongoURI, dbName)

	// Gin Setup
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	// Middleware chain
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register all API routes
	routes.RegisterRoutes(r, db)

	// Start Server
	log.Printf("🚀 Gin server running on http://localhost:%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
