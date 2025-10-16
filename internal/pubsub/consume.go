package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %v", err)
	}

	const noWait = false

	durable := false
	if queueType == SimpleQueueDurable {
		durable = true
	}

	autoDelete := false
	exclusive := false
	if queueType == SimpleQueueTransient {
		autoDelete = true
		exclusive = true
	}

	queue, err := channel.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create queue: %v", err)
	}

	err = channel.QueueBind(queueName, key, exchange, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %v", err)
	}

	return channel, queue, nil

}
