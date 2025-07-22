package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connection_strign := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connection_strign)
	if err != nil {
		fmt.Println("Error connection: ", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connection Sucessful...")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	fmt.Println("\nReceived signal", sig)
	fmt.Println("Shutting down program")
}
