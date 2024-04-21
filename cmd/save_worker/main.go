package main

import (
	"context"
	"log"
	"time"

	"github.com/kordyd/go-crawler/internal/db/mongodb"
	"github.com/kordyd/go-crawler/internal/services"
	"github.com/redis/go-redis/v9"
)

func main() {

	redisClient := redis.NewClient((&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}))

	initialCount, err := redisClient.DBSize(context.Background()).Result()

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

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	coll := client.Database(mongoDBName).Collection(mongoCollectionName)

	for {
		<-ticker.C
		currentCount, err := redisClient.DBSize(context.Background()).Result()

		if err != nil {
			log.Panicln(err)
		}

		if currentCount == initialCount {
			log.Println("Equal")

			if currentCount != 0 {
				services.SaveParsedUrls(coll, redisClient)
			}

		}

		if currentCount >= 200 {
			log.Println("Download content to db")
			services.SaveParsedUrls(coll, redisClient)
			log.Println("Download data complete")
		}

		initialCount = currentCount

	}

}
