package main

import (
	"context"
	"log"
	"time"

	"github.com/kordyd/go-crawler/internal/db/mongodb"
	"github.com/kordyd/go-crawler/internal/db/redis"
	"github.com/kordyd/go-crawler/internal/services"
)

func main() {
	redisClient := redis.Connect()

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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		currentCount, err := redisClient.DBSize(context.Background()).Result()

		if err != nil {
			log.Panicln(err)
		}

		if currentCount == initialCount {
			log.Println("Equal")
			coll := client.Database(mongoDBName).Collection(mongoCollectionName)

			if currentCount != 0 {
				services.SaveParsedUrls(coll, redisClient)
			}

		} else {
			log.Println("Not equal")
		}

		initialCount = currentCount

	}

}
