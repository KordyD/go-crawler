package rabbitmq

import (
	errorhandlers "github.com/kordyd/go-crawler/internal/error_handlers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, func() error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	errorhandlers.FailOnError(err)

	close := func() error {
		err := conn.Close()
		return err
	}

	return conn, close
}
