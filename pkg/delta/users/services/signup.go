package services

import (
	"context"
	"errors"
	"time"

	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *database.Database
}

func NewUserService(db *database.Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) SignUp(ctx context.Context, username, email, password string) (*models.User, error) {
	coll := s.db.Mongo.Database(s.db.DBName).Collection("users")
	count, err := coll.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Email:     email,
		Password:  string(hashedBytes),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
