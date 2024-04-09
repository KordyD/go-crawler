package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Url struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Url    string             `bson:"url"`
	Parsed bool               `bson:"parsed"`
}

func GetNotParsedUrls() []Url {
	mongoDBName := "Crawler"
	mongoCollectionName := "URLs"

	client, disconnect, err := MongodbConnect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	coll := client.Database(mongoDBName).Collection(mongoCollectionName)

	filter := bson.M{"parsed": false}
	cursor, err := coll.Find(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var urls []Url
	if err := cursor.All(context.TODO(), &urls); err != nil {
		log.Fatal(err)
	}

	return urls

}
