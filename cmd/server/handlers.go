package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

func handlerLogs() func(gamelog routing.GameLog) pubsub.Acktype {
	return func(gamelog routing.GameLog) pubsub.Acktype {
		defer fmt.Print("> ")

		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Printf("error writing log: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
