package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connection_strign := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connection_strign)
	if err != nil {
		fmt.Println("Error getting connection: ", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection Sucessful...")

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println("Error creating channel:", err)
		return
	}
	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		if input[0] == "pause" {
			fmt.Println("Sending a pause message")
			err = pubsub.PublishJSON(
				channel, routing.ExchangePerilDirect, routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				fmt.Println("Error publishing:", err)
				return
			}
			fmt.Println("Sucessfully Published pause message")
			continue
		}
		if input[0] == "resume" {
			fmt.Println("Sending a resume message")
			err = pubsub.PublishJSON(
				channel, routing.ExchangePerilDirect, routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				fmt.Println("Error publishing:", err)
				return
			}
			fmt.Println("Sucessfully Published resume message")
			continue
		}
		if input[0] == "quit" {
			fmt.Println("Exiting..")
			break
		}
		fmt.Println("Couldnt understand:", input[0])
	}

	fmt.Println("Shutting down program")
}
