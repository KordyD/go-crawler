package main

import (
	"context"
	"log"
	"sync"

	"github.com/kordyd/go-crawler/internal/entities"
	rabbitmq "github.com/kordyd/go-crawler/internal/rabbitMQ"
	"github.com/kordyd/go-crawler/internal/scrapper"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
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

	err = ch.Qos(
		1,
		0,
		false,
	)

	if err != nil {
		log.Panicln(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicln(err)
	}

	contentRedis := redis.NewClient((&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}))

	defer contentRedis.Close()

	linksToParse := redis.NewClient((&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}))

	defer linksToParse.Close()

	var wg sync.WaitGroup

	numberOfThreads := 5

	parsedData := make(chan entities.Url, numberOfThreads)
	parsedUrls := make(chan string, numberOfThreads)

	doneChan := make(chan string, 100)

	for i := 0; i < numberOfThreads; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for data := range parsedData {
				if data.Error != "" {
					_, err = linksToParse.SAdd(context.Background(), "errors", data.Link).Result()
					if err != nil {
						log.Println(err)
					}
					log.Println("Error with url", data.Error)
					// doneChan <- entities.Url{Link: data.Link, Parsed: false, Error: data.Error}
					doneChan <- data.Link
					continue
				}
				_, err = contentRedis.Set(context.Background(), data.Link, data.Content, 0).Result()
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("Set data in redis: %s", data.Link)
				_, err = linksToParse.SAdd(context.Background(), "parsed", data.Link).Result()
				if err != nil {
					log.Println(err)
				}
				// doneChan <- entities.Url{Link: data.Link, Parsed: true}
				doneChan <- data.Link
			}
		}()
		go func() {
			defer wg.Done()
			for url := range parsedUrls {
				isExist, err := linksToParse.SIsMember(context.Background(), "parsed", url).Result()
				if err != nil {
					log.Println(err)
					continue
				}
				if isExist {
					log.Println(url, "already parsed")
					continue
				}
				isAdded, err := linksToParse.SAdd(context.Background(), "toParse", url).Result()
				if err != nil {
					log.Println(err)
					continue
				}
				if isAdded == 1 {
					log.Printf("Set url to parse in redis: %s", url)
				} else {
					log.Printf("Value already exists %s", url)
				}
			}
		}()
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)

		go scrapper.Scrapper(string(d.Body), parsedData, parsedUrls)
		go func(d amqp091.Delivery) {
			message := <-doneChan
			err := ch.PublishWithContext(context.Background(),
				"",
				d.ReplyTo,
				false,
				false,
				amqp091.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(message),
				})
			if err != nil {
				log.Println(err)
			}
			log.Println("Done!!!!", message)
			d.Ack(false)
		}(d)

	}

	wg.Wait()

}
