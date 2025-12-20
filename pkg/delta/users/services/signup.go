package services

import (
	"context"
	"errors"
	"time"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// SignUp handles the business logic for registering a new user
func (s *UserService) SignUp(ctx context.Context, username, email, password string) (*models.User, error) {
	coll := s.db.Mongo.Database(s.db.DBName).Collection("users")

	// 1. Check if user exists
	count, err := coll.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	// 2. Hash Password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create User Model
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Email:     email,
		Password:  string(hashedBytes),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// 4. Persist
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
