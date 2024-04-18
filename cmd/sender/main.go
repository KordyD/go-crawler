package main

import (
	"context"
	"log"
	"math/rand"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urls := services.GetNotParsedUrls()

	for _, url := range urls {
		corrId := randomString(32)
		body := url.Link
		err := ch.PublishWithContext(ctx,
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

		for d := range msgs {
			if corrId == d.CorrelationId {
				log.Println(string(d.Body))
				break
			}
		}
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

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
