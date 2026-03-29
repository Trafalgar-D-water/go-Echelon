package mongodb

import (
	"context"
	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SessionStore struct {
	Collection *mongo.Collection
}

func (s *SessionStore) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	_, err := s.Collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionStore) GetSessionByUserID(ctx context.Context, userID string) (*models.Session, error) {
	var session models.Session
	err := s.Collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionStore) UpdateSession(ctx context.Context, sessionID string, newRefreshToken string, expiresAt time.Time) error {
	oid, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return err
	}
	_, err = s.Collection.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{
			"$set": bson.M{
				"refreshToken": newRefreshToken,
				"expiresAt":    expiresAt,
			},
		})
	return err
}
