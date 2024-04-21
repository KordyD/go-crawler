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
	keys := redisClient.Keys(context.Background(), "*").Val()
	for _, key := range keys {
		// value := redisClient.HGetAll(context.Background(), key).Val()
		value, err := redisClient.Get(context.Background(), key).Result()
		if err != nil {
			log.Println(err)
		}
		_, err = redisClient.Del(context.Background(), key).Result()
		if err != nil {
			log.Println(err)
		}

		// url := entities.Url{Link: value["link"], Parsed: value["parsed"] != "0", Error: value["error"], Content: value["content"]}
		url := entities.Url{Link: key, Content: value}

		// options := options.Replace().SetUpsert(true)
		// filter := bson.M{"link": url.Link}
		// _, err = mongoCollection.ReplaceOne(context.Background(), filter, url, options)
		// if err != nil {
		// 	panic(err)
		// }

		_, err = mongoCollection.InsertOne(context.Background(), url)
		if err != nil {
			log.Printf("Failed to insert data for key %s: %v", key, err)
		}
	}

	// fmt.Println("Data migration from Redis to MongoDB complete.")

}
