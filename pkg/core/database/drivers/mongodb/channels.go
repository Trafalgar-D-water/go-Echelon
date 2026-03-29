package mongodb

import (
	"context"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelStore struct {
	Collection *mongo.Collection
}

func (s *ChannelStore) CreateChannel(ctx context.Context, channel *models.Channel) error {
	_, err := s.Collection.InsertOne(ctx, channel)
	return err
}

func (s *ChannelStore) GetChannelByID(ctx context.Context, id string) (*models.Channel, error) {
	var channel models.Channel
	err := s.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (s *ChannelStore) DeleteChannel(ctx context.Context, id string) error {
	_, err := s.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
