package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/kordyd/go-crawler/db"
	"github.com/kordyd/go-crawler/scrapper"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"url_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	client := db.Connect()
	ctx := context.Background()

	// var forever chan struct{}

	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received a message: %s", d.Body)

	// 		var wg sync.WaitGroup

	// 		parsedUrls := make(chan string, 5)
	// 		fetchedBody := make(chan string, 5)

	// 		go func() {
	// 			wg.Wait()
	// 			close(parsedUrls)
	// 			close(fetchedBody)
	// 		}()

	// 		wg.Add(1)
	// 		go scrapper.Scrapper(string(d.Body), parsedUrls, fetchedBody, &wg)

	// 		for url := range parsedUrls {
	// 			client.Set(ctx, url, 1, 0).Result()
	// 			fmt.Println(url)
	// 		}

	// 		// for body := range fetchedBody {
	// 		// 	fmt.Println(body)
	// 		// }

	// 		// dotCount := bytes.Count(d.Body, []byte("."))
	// 		// t := time.Duration(dotCount)
	// 		// time.Sleep(t * time.Second)
	// 		log.Printf("Done")
	// 	}
	// }()

	var wg sync.WaitGroup

	parsedUrls := make(chan string, 5)
	fetchedBody := make(chan string, 5)

	for i := 0; i < 5; i++ { // Change 5 to the desired number of concurrent scrapers
		wg.Add(1)
		go func() {
			defer wg.Done()

			for url := range parsedUrls {
				client.Set(ctx, url, 1, 0).Result()
				fmt.Println(url)
			}
		}()
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		go func(msg []byte) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from panic:", r)
				}
			}()
			scrapper.Scrapper(string(msg), parsedUrls, fetchedBody)
		}(d.Body)
	}

	wg.Wait()

	// <-forever
}
