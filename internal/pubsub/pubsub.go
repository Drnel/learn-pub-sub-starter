package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	json_bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(
		context.Background(), exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        json_bytes,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

type simpleQueueType int

const (
	Durable simpleQueueType = iota
	Transient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange, queueName, key string,
	queueType simpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	queue, err := channel.QueueDeclare(
		queueName,
		queueType == Durable,
		queueType == Transient,
		queueType == Transient,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return channel, queue, nil
}
