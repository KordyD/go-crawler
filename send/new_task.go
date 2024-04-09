package main

import (
	"context"
	"log"
	"time"

	errorhandler "github.com/kordyd/go-crawler/error_handler"
	"github.com/kordyd/go-crawler/mongodb"
	rabbitmq "github.com/kordyd/go-crawler/rabbitMQ"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn := rabbitmq.Connect()
	defer conn.Close()

	ch, err := conn.Channel()
	errorhandler.FailOnError(err, "Failed to pen a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"url_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	errorhandler.FailOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	urls := mongodb.GetNotParsedUrls()

	for _, url := range urls {
		sendMsg(url.Url, ch, q, ctx)
	}

	// jsonData, err := json.MarshalIndent(urls, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", jsonData)

	// // Print the URLs
	// for _, url := range urls {
	// 	fmt.Println(url)
	// }
	// conn := rabbitmq.Connect()
	// defer conn.Close()

	// ch, err := conn.Channel()
	// errorhandler.FailOnError(err, "Failed to open a channel")
	// defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"url_queue", // name
	// 	false,       // durable
	// 	false,       // delete when unused
	// 	false,       // exclusive
	// 	false,       // no-wait
	// 	nil,         // arguments
	// )
	// errorhandler.FailOnError(err, "Failed to declare a queue")
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// body := os.Args[1]

	// sendMsg(body, ch, q, ctx)

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
	errorhandler.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n", body)
}
