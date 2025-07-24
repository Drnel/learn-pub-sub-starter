package main

import (
	"fmt"
	"os"
	"os/signal"

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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	fmt.Println("\nReceived signal", sig)
	fmt.Println("Shutting down program")
}
