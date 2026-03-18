package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/delta/middleware"
	"github.com/go-Echelon/go-Echelon/pkg/delta/routes"

	"github.com/go-Echelon/go-Echelon/pkg/core/config"

	// Swagger API Docs
	_ "github.com/go-Echelon/go-Echelon/api"
)

// @title           Go-Echelon API (Delta)
// @version         1.0
// @description     This is the REST API server for the Go-Echelon backend.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@go-echelon.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	log.Println("🔌 Connecting to MongoDB...")

	db, err := database.Connect(cfg.MongoURI, cfg.DBName)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	log.Printf("✅ MongoDB connected: %s/%s\n", cfg.MongoURI, cfg.DBName)

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
	log.Printf("🚀 Gin server running on http://localhost:%s", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
