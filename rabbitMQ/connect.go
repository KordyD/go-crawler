package rabbitmq

import (
	errorhandler "github.com/kordyd/go-crawler/error_handler"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	errorhandler.FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}
