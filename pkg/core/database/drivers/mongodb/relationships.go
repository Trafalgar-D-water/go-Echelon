package mongodb

import (
	"context"
	"time"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RelationshipStore struct {
	Collection *mongo.Collection
}

func (s *RelationshipStore) SendRequest(ctx context.Context, actorID, targetID string) error {
	now := time.Now()
	outgoing := models.Relationship{
		ActorID: actorID, TargetID: targetID, Type: models.RelOutgoing, CreatedAt: now, UpdatedAt: now,
	}
	incoming := models.Relationship{
		ActorID: targetID, TargetID: actorID, Type: models.RelIncoming,
		CreatedAt: now, UpdatedAt: now,
	}

	_, err := s.Collection.InsertMany(ctx, []interface{}{
		outgoing,
		incoming,
	})

	return err
}

func (s *RelationshipStore) AcceptRequest(ctx context.Context, actorID, targetID string) error {
	filter := bson.M{
		"$or": []bson.M{
			{"actor_id": actorID, "target_id": targetID},
			{"actor_id": targetID, "target_id": actorID},
		},
	}
	// Update BOTH sides of the relationship to 'RelFriend' (Type 1)
	update := bson.M{
		"$set": bson.M{
			"type":       models.RelFriend,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err := s.Collection.UpdateMany(ctx, filter, update)
	return err
}

func (s *RelationshipStore) GetUserRelationships(ctx context.Context, userID string) ([]*models.Relationship, error) {
	cursor, err := s.Collection.Find(ctx, bson.M{"actor_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rels []*models.Relationship
	if err := cursor.All(ctx, &rels); err != nil {
		return nil, err
	}
	return rels, nil
}
