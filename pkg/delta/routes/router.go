package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/delta/routes/users"
)

// RegisterRoutes registers all API route groups on the engine.
func RegisterRoutes(r *gin.Engine, db *database.Database) {
	api := r.Group("/api/v1")

	// User routes
	users.RegisterRoutes(api, db)
}
