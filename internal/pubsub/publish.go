package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {

	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(val)
	if err != nil {
		fmt.Println("couldn't encode value")
		return err
	}

	const mandatory = false
	const inmediate = false
	err = ch.PublishWithContext(context.Background(), exchange, key, mandatory, inmediate, amqp.Publishing{
		ContentType: "application/gob",
		Body:        data.Bytes(),
	})

	return err
}

func PublishGameLog(publishCh *amqp.Channel, username, msg string) error {
	return PublishGob(
		publishCh,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			Username:    username,
			CurrentTime: time.Now(),
			Message:     msg,
		},
	)
}
