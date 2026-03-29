package drivers

import (
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
	"github.com/go-Echelon/go-Echelon/pkg/core/database/drivers/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database

	userStore    database.UserStore
	sessionStore database.SessionStore
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
	}
}

func (m *MongoDB) Users() database.UserStore {
	return m.userStore
}

func (m *MongoDB) Sessions() database.SessionStore {
	return m.sessionStore
}
