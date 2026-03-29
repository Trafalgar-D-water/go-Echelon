package mongodb

import (
	"context"

	"time"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	Collection *mongo.Collection
}

func (s *UserStore) CreateUser(ctx context.Context, user *models.User) error {
	_, err := s.Collection.InsertOne(ctx, user)
	return err
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) CountByEmail(ctx context.Context, email string) (int64, error) {
	return s.Collection.CountDocuments(ctx, bson.M{"email": email})
}

func (s *UserStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.Collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) VerifyUserOTP(ctx context.Context, email string, otp string) (bool, error) {
	filter := bson.M{
		"email": email,
		"otp":   otp,
	}

	update := bson.M{
		"$set": bson.M{
			"is_verified": true,
			"updated_at":  time.Now().UTC(),
		},
		"$unset": bson.M{
			"otp": "",
		},
	}

	result, err := s.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return result.MatchedCount > 0, nil
}
