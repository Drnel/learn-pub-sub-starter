package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connection_strign := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connection_strign)
	if err != nil {
		fmt.Println("Error getting connection: ", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection Sucessful...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Error getting username:", err)
		return
	}
	channel, queue, err := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		fmt.Println("Error getting channel and queue:", err)
		return
	}
	_, _ = channel, queue
	game_state := gamelogic.NewGameState(username)
	pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		"pause."+username,
		routing.PauseKey,
		pubsub.Transient,
		pubsub.HandlerPause(game_state),
	)
Loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			err = game_state.CommandSpawn(input)
			if err != nil {
				fmt.Println("Couldnt parse spawn command:", err)
				continue Loop
			}
		case "move":
			_, err = game_state.CommandMove(input)
			if err == nil {
				fmt.Println("Sucessfully Moved")
			} else {
				fmt.Println("Error, moving:", err)
			}
		case "status":
			game_state.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			fmt.Println("Exiting...")
			break Loop
		default:
			fmt.Println("Could'nt understand command:", input[0])
		}
	}

	fmt.Println("Shutting down program")
}
