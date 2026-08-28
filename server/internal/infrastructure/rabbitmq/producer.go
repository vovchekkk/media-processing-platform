package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	connectionManager *ConnectionManager
	queueName         string
	log               *slog.Logger
}

func NewProducer(connectionManager *ConnectionManager, queueName string, log *slog.Logger) *Producer {
	return &Producer{
		connectionManager: connectionManager,
		queueName:         queueName,
		log:               log,
	}
}

func (producer *Producer) Publish(ctx context.Context, body []byte) error {
	conn, err := producer.connectionManager.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get active rabbitmq connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open rabbitmq channel: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		producer.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare rabbitmq queue: %w", err)
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		producer.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/x-protobuf",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %w", producer.queueName, err)
	}

	producer.log.Debug("task successfully published to rabbitmq", slog.String("queue", producer.queueName))
	return nil
}
