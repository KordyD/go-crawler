package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kordyd/go-crawler/internal/db/mongodb"
	"github.com/kordyd/go-crawler/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	mongoDBName := "Crawler"
	mongoCollectionName := "URLs"

	client, disconnect, err := mongodb.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	coll := client.Database(mongoDBName).Collection(mongoCollectionName)

	insertTestData(coll)

}

func insertTestData(mongoCollection *mongo.Collection) {

	_, err := mongoCollection.DeleteMany(context.TODO(), bson.D{})

	if err != nil {
		log.Fatalln("Failed to delete data")
	}

	urls := []interface{}{
		entities.Url{Link: "https://redis.uptrace.dev/", Parsed: false, Error: ""},
		entities.Url{Link: "https://pkg.go.dev/", Parsed: false, Error: ""},
		entities.Url{Link: "https://www.rabbitmq.com/", Parsed: false, Error: "Smth error"},
		entities.Url{Link: "https://habr.com/", Parsed: true, Error: ""},
		entities.Url{Link: "https://stackoverflow.com/", Parsed: false, Error: ""},
	}

	_, err = mongoCollection.InsertMany(context.TODO(), urls)

	if err != nil {
		log.Fatalln("Failed to insert data")
	}

	fmt.Println("Data insert complete.")

}
