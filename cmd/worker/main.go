package main

import (
	"context"
	"log"
	"sync"

	"github.com/kordyd/go-crawler/internal/db/redis"
	"github.com/kordyd/go-crawler/internal/entities"
	rabbitmq "github.com/kordyd/go-crawler/internal/rabbitMQ"
	"github.com/kordyd/go-crawler/internal/scrapper"
)

func main() {

	conn, close := rabbitmq.Connect()
	defer close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicln(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"url_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Panicln(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicln(err)
	}

	client := redis.Connect()

	var wg sync.WaitGroup

	numberOfThreads := 5

	parsedData := make(chan entities.Url, numberOfThreads)
	parsedUrls := make(chan string, numberOfThreads)

	for i := 0; i < numberOfThreads; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for data := range parsedData {
				_, err = client.HSet(context.TODO(), data.Link, entities.Url{
					Link:    data.Link,
					Parsed:  true,
					Error:   data.Error,
					Content: data.Content,
				}).Result()
				if err != nil {
					log.Println(err, "err in parseddata")
					continue
				}
				log.Printf("Set data in redis: %s", data.Link)
			}
		}()
		go func() {
			defer wg.Done()
			for url := range parsedUrls {
				exists, err := client.Exists(context.TODO(), url).Result()
				if err != nil {
					log.Println(err)
					continue
				}
				if exists == 0 {
					_, err = client.HSet(context.TODO(), url, entities.Url{
						Link:    url,
						Parsed:  false,
						Error:   "",
						Content: ""}).Result()
					if err != nil {
						log.Println(err, "err in parsedurl")
						continue
					}
					log.Printf("Set url to parse in redis: %s", url)
				} else {
					log.Printf("Key already exists in Redis: %s", url)
				}
			}
		}()
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		go scrapper.Scrapper(string(d.Body), parsedData, parsedUrls)
	}

	wg.Wait()

}
