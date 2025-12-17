package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Mongo  *mongo.Client
	DBName string
}

func Connect(mongoUri, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myOpts := options.Client().ApplyURI(mongoUri)
	mClient, err := mongo.Connect(ctx, myOpts)

	if err != nil {
		return nil, err
	}

	if err := mClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Database{
		Mongo:  mClient,
		DBName: dbName,
	}, nil

}
