package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}

	defer connection.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("could not create channel in RabbitMQ: %v", err)
	}

	err = pubsub.SubscribeGob(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerLogs(),
	)
	if err != nil {
		log.Fatalf("could not starting consuming logs: %v", err)
	}

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()

		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			fmt.Println("Publishing pause message...")
			err = updatePlayingState(channel, true)
			if err != nil {
				log.Fatalf("could not handle pause: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resume message...")
			err = updatePlayingState(channel, false)
			if err != nil {
				log.Fatalf("could not handle pause: %v", err)
			}
		case "quit":
			fmt.Println("Closing server...")
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}
