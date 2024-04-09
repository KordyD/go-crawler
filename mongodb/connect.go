package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongodbConnect() (*mongo.Client, func() error, error) {

	mongoURI := "mongodb://localhost:27017"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))

	if err != nil {
		return nil, nil, err
	}

	closeFunc := func() error {
		err := client.Disconnect(context.TODO())
		return err
	}

	return client, closeFunc, nil

}
