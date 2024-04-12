package main

import (
	"context"
	"log"
	"time"

	errorhandlers "github.com/kordyd/go-crawler/internal/error_handlers"
	rabbitmq "github.com/kordyd/go-crawler/internal/rabbitMQ"
	"github.com/kordyd/go-crawler/internal/services"
	amqp "github.com/rabbitmq/amqp091-go"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urls := services.GetNotParsedUrls()

	for _, url := range urls {
		sendMsg(url.Link, ch, q, ctx)
	}
}

func sendMsg(body string, ch *amqp.Channel, q amqp.Queue, ctx context.Context) {

	err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	errorhandlers.FailOnError(err)

	log.Printf(" [x] Sent %s\n", body)
}
