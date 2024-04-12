package services

import (
	"context"

	"github.com/kordyd/go-crawler/internal/db/mongodb"
	"github.com/kordyd/go-crawler/internal/entities"
	errorhandlers "github.com/kordyd/go-crawler/internal/error_handlers"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNotParsedUrls() []entities.Url {
	mongoDBName := "Crawler"
	mongoCollectionName := "URLs"

	client, disconnect, err := mongodb.Connect()

	errorhandlers.FailOnError(err)

	defer func() {
		err := disconnect()
		errorhandlers.FailOnError(err)
	}()

	coll := client.Database(mongoDBName).Collection(mongoCollectionName)

	filter := bson.M{"parsed": false}
	cursor, err := coll.Find(context.TODO(), filter)

	errorhandlers.FailOnError(err)

	defer cursor.Close(context.TODO())

	var urls []entities.Url
	err = cursor.All(context.TODO(), &urls)

	errorhandlers.FailOnError(err)

	return urls

}
