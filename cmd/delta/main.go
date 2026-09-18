package main

import (
	"context"
	"log/slog"
	"time"
	// Swagger API Docs

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/go-Echelon/go-Echelon/api"
	"github.com/go-Echelon/go-Echelon/pkg/core/config"
	"github.com/go-Echelon/go-Echelon/pkg/core/database/drivers"
	"github.com/go-Echelon/go-Echelon/pkg/core/observability"
	"github.com/go-Echelon/go-Echelon/pkg/delta/middleware"
	"github.com/go-Echelon/go-Echelon/pkg/delta/routes"
	"github.com/go-Echelon/go-Echelon/pkg/delta/util"
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

	logger := observability.NewLogger(observability.Config{
		Level:   cfg.LogLevel,
		Format:  cfg.LogFormat,
		Service: "delta",
	})

	slog.SetDefault(logger)
	gin.SetMode(cfg.GinMode)

	if err := util.ConfigureJWT(
		cfg.AccessTokenSecret,
		cfg.RefreshTokenSecret,
	); err != nil {
		logger.With(
			"component", "delta.bootstrap",
		).Error(
			"invalid JWT configuration",
			"error", err,
		)
		return
	}

	logger.With(
		"component", "delta.bootstrap",
	).Info("connecting to mongodb")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.With(
			"component", "delta.bootstrap",
		).Error(
			"failed to connect to mongodb",
			"error", err,
		)
		return
	}
	defer client.Disconnect(context.Background())

	// Ping the DB to ensure connection is actually successful
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.With(
			"component", "delta.bootstrap",
		).Error(
			"failed to ping mongodb",
			"error", err,
		)
		return
	}

	logger.With(
		"component", "delta.bootstrap",
	).Info(
		"mongodb connected",
		"database", cfg.DBName,
	)

	db := drivers.New(client, cfg.DBName)
	// Gin Setup
	gin.SetMode(gin.DebugMode)
	r := gin.New()

	// Middleware chain
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register all API routes
	routes.RegisterRoutes(r, db)

	// Start Server
	logger.With(
		"component", "delta.bootstrap",
	).Info(
		"delta server started",
		"port", cfg.Port,
	)

	if err := r.Run(":" + cfg.Port); err != nil {
		logger.With(
			"component", "delta.bootstrap",
		).Error(
			"delta server stopped unexpectedly",
			"error", err,
		)
	}
}
