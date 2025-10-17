package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func updatePlayingState(ch *amqp.Channel, isPaused bool) error {

	playingState := routing.PlayingState{
		IsPaused: isPaused,
	}

	err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, playingState)
	if err != nil {
		return fmt.Errorf("could not publish to RabbitMQ: %v", err)
	}

	return nil
}
