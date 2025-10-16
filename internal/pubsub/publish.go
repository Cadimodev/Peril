package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	data, err := json.Marshal(val)
	if err != nil {
		fmt.Println("couldn't marshal value")
		return err
	}

	const mandatory = false
	const inmediate = false
	err = ch.PublishWithContext(context.Background(), exchange, key, mandatory, inmediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})

	return err
}
