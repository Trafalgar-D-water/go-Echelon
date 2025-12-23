package services

import (
	"context"
	"errors"

	"github.com/go-Echelon/go-Echelon/pkg/core/config"
	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	coll := s.db.Mongo.Database(s.db.DBName).Collection("users")

	var user models.User
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	token, err := config.GenerateToken(user.ID.Hex())
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}
