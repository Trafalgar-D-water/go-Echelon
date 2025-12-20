package database

import "go.mongodb.org/mongo-driver/mongo"

const (
	UsersCollection    = "users"
	ChannelsCollection = "channels"
	MessagesCollection = "messages"
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
