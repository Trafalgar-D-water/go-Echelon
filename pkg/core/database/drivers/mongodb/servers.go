package mongodb

import (
	"context"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerStore struct {
	Collection *mongo.Collection
}

func (s *ServerStore) CreateServer(ctx context.Context, server *models.Server) error {
	_, err := s.Collection.InsertOne(ctx, server)
	return err
}

func (s *ServerStore) GetServerByID(ctx context.Context, id string) (*models.Server, error) {
	var server models.Server
	err := s.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&server)
	if err != nil {
		return nil, err
	}
	return &server, nil
}
