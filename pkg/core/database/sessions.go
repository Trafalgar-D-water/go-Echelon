package database

import (
	"context"
	"time"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
)

type SessionStore interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
	GetSessionByUserID(ctx context.Context, userID string) (*models.Session, error)
	UpdateSession(ctx context.Context, sessionID string, newRefreshToken string, expiresAt time.Time) error
}
