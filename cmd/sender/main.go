package main

import (
	"context"
	"log"
	"sync"
	"time"

	uuid "github.com/google/uuid"
	errorhandlers "github.com/kordyd/go-crawler/internal/error_handlers"
	rabbitmq "github.com/kordyd/go-crawler/internal/rabbitMQ"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var (
	repliesMutex    sync.Mutex
	receivedReplies = make(map[string]bool)
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
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
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

	linksToParse := redis.NewClient((&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}))

	defer linksToParse.Close()

	go consumeReplies(msgs, linksToParse)

	for {
		log.Println("Start timer")
		<-time.After(5 * time.Second)

		log.Println("Fetching urls")
		urls, err := linksToParse.SMembers(context.Background(), "toParse").Result()
		if err != nil {
			log.Panicln(err)
		}

		for _, url := range urls {

			corrId := uuid.NewString()
			body := url
			err := ch.PublishWithContext(context.Background(),
				"",          // exchange
				"url_queue", // routing key
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ReplyTo:       q.Name,
					ContentType:   "text/plain",
					CorrelationId: corrId,
					Body:          []byte(body),
				})
			errorhandlers.FailOnError(err)

			log.Printf(" [x] Sent %s\n", body)
			trackSentMessage(corrId)

		}
		waitForReplies(len(urls))

	}

}

func consumeReplies(msgs <-chan amqp.Delivery, redisClient *redis.Client) {
	for d := range msgs {
		log.Println(string(d.Body))
		_, err := redisClient.SRem(context.Background(), "toParse", string(d.Body)).Result()
		if err != nil {
			log.Println(err)
		}
		trackReceivedMessage(d.CorrelationId)
	}
}

func trackSentMessage(corrID string) {
	repliesMutex.Lock()
	defer repliesMutex.Unlock()
	receivedReplies[corrID] = false
}

func trackReceivedMessage(corrID string) {
	repliesMutex.Lock()
	defer repliesMutex.Unlock()
	receivedReplies[corrID] = true
}

func waitForReplies(expectedReplies int) {
	for {
		repliesMutex.Lock()
		receivedCount := 0
		for _, received := range receivedReplies {
			if received {
				receivedCount++
			}
		}
		repliesMutex.Unlock()

		if receivedCount >= expectedReplies {
			log.Println("All replies received")
			break
		}

		time.Sleep(1 * time.Second)
	}

	// Clear received replies for the next iteration
	repliesMutex.Lock()
	receivedReplies = make(map[string]bool)
	repliesMutex.Unlock()

}
