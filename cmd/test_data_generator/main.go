package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PuerkitoBio/purell"
	"github.com/kordyd/go-crawler/internal/db/mongodb"
	"github.com/kordyd/go-crawler/internal/entities"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	redisClient := redis.NewClient((&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}))

	defer redisClient.Close()

	links := []entities.Url{
		{Link: "https://redis.uptrace.dev/", Parsed: false, Error: ""},
		{Link: "https://pkg.go.dev/", Parsed: false, Error: ""},
		{Link: "https://www.rabbitmq.com/", Parsed: false, Error: "Smth error"},
		{Link: "https://habr.com/", Parsed: true, Error: ""},
		{Link: "https://stackoverflow.com/", Parsed: false, Error: ""},
		{Link: "https://medium.com/", Parsed: false, Error: ""},
	}

	var urls []string

	flags := purell.FlagsUsuallySafeGreedy | purell.FlagRemoveDirectoryIndex | purell.FlagRemoveFragment | purell.FlagRemoveDuplicateSlashes | purell.FlagRemoveWWW | purell.FlagSortQuery

	for _, link := range links {
		normilizedUrl, err := purell.NormalizeURLString(link.Link, flags)
		if err != nil {
			log.Panicln(err)
		}
		urls = append(urls, normilizedUrl)
	}

	_, err := redisClient.FlushAll(context.Background()).Result()
	if err != nil {
		log.Panicln(err)
	}

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

	_, err = coll.DeleteMany(context.Background(), bson.D{})

	if err != nil {
		log.Fatalln("Failed to delete data")
	}

	number, err := redisClient.SAdd(context.Background(), "toParse", urls).Result()

	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("Data insert complete:", number)

	// mongoDBName := "Crawler"
	// mongoCollectionName := "URLs"

	// client, disconnect, err := mongodb.Connect()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer func() {
	// 	if err := disconnect(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// coll := client.Database(mongoDBName).Collection(mongoCollectionName)

	// insertTestData(coll)

}

// func insertTestData(mongoCollection *mongo.Collection) {

// 	_, err := mongoCollection.DeleteMany(context.Background(), bson.D{})

// 	if err != nil {
// 		log.Fatalln("Failed to delete data")
// 	}

// 	_, err = mongoCollection.Indexes().CreateOne(
// 		context.Background(),
// 		mongo.IndexModel{
// 			Keys:    bson.M{"link": 1},
// 			Options: options.Index().SetUnique(true),
// 		},
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// https://go-colly.org/

// 	links := []entities.Url{
// 		{Link: "https://redis.uptrace.dev/", Parsed: false, Error: ""},
// 		{Link: "https://pkg.go.dev/", Parsed: false, Error: ""},
// 		{Link: "https://www.rabbitmq.com/", Parsed: false, Error: "Smth error"},
// 		{Link: "https://habr.com/", Parsed: true, Error: ""},
// 		{Link: "https://stackoverflow.com/", Parsed: false, Error: ""},
// 		{Link: "https://medium.com/", Parsed: false, Error: ""},
// 	}

// 	var urls []interface{}

// 	flags := purell.FlagsUsuallySafeGreedy | purell.FlagRemoveDirectoryIndex | purell.FlagRemoveFragment | purell.FlagRemoveDuplicateSlashes | purell.FlagRemoveWWW | purell.FlagSortQuery

// 	for _, link := range links {
// 		normilizedUrl, err := purell.NormalizeURLString(link.Link, flags)
// 		if err != nil {
// 			log.Panicln(err)
// 		}
// 		link.Link = normilizedUrl
// 		urls = append(urls, link)
// 	}

// 	_, err = mongoCollection.InsertMany(context.Background(), urls)

// 	if err != nil {
// 		log.Fatalln("Failed to insert data")
// 	}

// 	fmt.Println("Data insert complete.")

// }
