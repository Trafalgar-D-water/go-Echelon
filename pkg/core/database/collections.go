package database

import "go.mongodb.org/mongo-driver/mongo"

const (
	UsersCollection    = "users"
	ChannelsCollection = "channels"
	MessagesCollection = "messages"
	SessionsCollection = "sessions"
)

func (d *Database) Users() *mongo.Collection {
	return d.Mongo.Database(d.DBName).Collection(UsersCollection)
}

func (d *Database) Channels() *mongo.Collection {
	return d.Mongo.Database(d.DBName).Collection(ChannelsCollection)
}

func (d *Database) Messages() *mongo.Collection {
	return d.Mongo.Database(d.DBName).Collection(MessagesCollection)
}

func (d *Database) Sessions() *mongo.Collection {
	return d.Mongo.Database(d.DBName).Collection(SessionsCollection)
}
