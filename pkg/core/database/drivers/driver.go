package drivers

import (
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/core/database/drivers/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database

	userStore         database.UserStore
	sessionStore      database.SessionStore
	relationshipStore database.RelationshipStore
}

func New(client *mongo.Client, dbName string) *MongoDB {
	db := client.Database(dbName)

	return &MongoDB{
		client: client,
		db:     db,

		userStore: &mongodb.UserStore{
			Collection: db.Collection("users"),
		},

		sessionStore: &mongodb.SessionStore{
			Collection: db.Collection("sessions"),
		},

		relationshipStore: &mongodb.RelationshipStore{
			Collection: db.Collection("relationships"),
		},
	}
}

func (m *MongoDB) Users() database.UserStore {
	return m.userStore
}

func (m *MongoDB) Sessions() database.SessionStore {
	return m.sessionStore
}

func (m *MongoDB) Relationships() database.RelationshipStore {
	return m.relationshipStore
}
