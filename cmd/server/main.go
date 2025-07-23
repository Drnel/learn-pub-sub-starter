package main

import (
	"fmt"
	"os"
	"os/signal"

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
	err = pubsub.PublishJSON(
		channel, routing.ExchangePerilDirect, routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)
	if err != nil {
		fmt.Println("Error publishing:", err)
		return
	}
	fmt.Println("Sucessfully Published")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	fmt.Println("\nReceived signal", sig)
	fmt.Println("Shutting down program")
}
