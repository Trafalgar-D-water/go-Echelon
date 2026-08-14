package database

import (
	"context"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
)

type RelationshipStore interface {
	SendRequest(ctx context.Context, actorID, targetID string) error
	AcceptRequest(ctx context.Context, actorID, targetID string) error
	GetUserRelationships(ctx context.Context, userID string) ([]*models.Relationship, error)
}
