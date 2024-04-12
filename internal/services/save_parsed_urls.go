package services

import (
	"context"
	"log"

	"github.com/kordyd/go-crawler/internal/entities"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveParsedUrls(mongoCollection *mongo.Collection, redisClient *redis.Client) {

	// Retrieve data from Redis
	keys := redisClient.Keys(context.TODO(), "*").Val()
	for _, key := range keys {
		value := redisClient.HGetAll(context.TODO(), key).Val()

		url := entities.Url{Link: value["link"], Parsed: value["parsed"] != "0", Error: value["error"]}

		// Here you can perform any transformation of data if necessary
		// For simplicity, we just directly store it in MongoDB
		_, err := mongoCollection.InsertOne(context.TODO(), url)
		if err != nil {
			log.Printf("Failed to insert data for key %s: %v", key, err)
		} else {
			log.Printf("Inserted data for key %s", key)
		}
	}

	// fmt.Println("Data migration from Redis to MongoDB complete.")

}
