package services

import (
	"context"

	"github.com/kordyd/go-crawler/internal/db/neo4jdb"
	"github.com/kordyd/go-crawler/internal/entities"
	errorhandlers "github.com/kordyd/go-crawler/internal/error_handlers"
	"go.neo4jdb.org/neo4j-driver/bson"
)

func GetNotParsedUrls() []entities.Url {
	neo4jDBName := "Crawler"
	neo4jCollectionName := "URLs"

	client, disconnect, err := neo4jdb.Connect()

	errorhandlers.FailOnError(err)

	defer func() {
		err := disconnect()
		errorhandlers.FailOnError(err)
	}()

	coll := client.Database(neo4jDBName).Collection(neo4jCollectionName)

	filter := bson.M{"parsed": false}
	cursor, err := coll.Find(context.Background(), filter)

	errorhandlers.FailOnError(err)

	defer cursor.Close(context.Background())

	var urls []entities.Url
	err = cursor.All(context.Background(), &urls)

	errorhandlers.FailOnError(err)

	return urls

}
