package rabbitmq

import (
	"context"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskHandler func(ctx context.Context, body []byte) error

type Consumer struct {
	connectionManager *ConnectionManager
	queueName         string
	log               *slog.Logger
}

func NewConsumer(connectionManager *ConnectionManager, queueName string, log *slog.Logger) *Consumer {
	return &Consumer{
		connectionManager: connectionManager,
		queueName:         queueName,
		log:               log,
	}
}

func (consumer *Consumer) StartListening(ctx context.Context, handler TaskHandler) error {
	for {
		select {
		case <-ctx.Done():
			consumer.log.Info("consumer worker context cancelled, stopping...")
			return nil
		default:
		}

		conn, err := consumer.connectionManager.GetConnection()
		if err != nil {
			consumer.log.Warn("waiting for active RabbitMQ connection...", slog.String("error", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			consumer.log.Error("failed to open channel for consumer", slog.String("error", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}

		_, err = ch.QueueDeclare(consumer.queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			consumer.log.Error("failed to declare queue", slog.String("error", err.Error()))
			_ = ch.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		msgs, err := ch.Consume(
			consumer.queueName,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			consumer.log.Error("failed to register consumer", slog.String("error", err.Error()))
			_ = ch.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		consumer.log.Info("consumer start listening for messages", slog.String("queue", consumer.queueName))

		for msg := range msgs {
			consumer.handleSingleMessage(ctx, msg, handler)
		}

		consumer.log.Warn("consumer channel closed, reconnecting consumer loop...")
		_ = ch.Close()
		time.Sleep(1 * time.Second)
	}
}

func (consumer *Consumer) handleSingleMessage(ctx context.Context, msg amqp.Delivery, handler TaskHandler) {
	consumer.log.Debug("received message from RabbitMQ", slog.Uint64("delivery_tag", msg.DeliveryTag))

	err := handler(ctx, msg.Body)
	if err != nil {
		consumer.log.Error("failed to handle message", slog.String("error", err.Error()))

		if err := msg.Nack(false, false); err != nil {
			consumer.log.Error("failed to nack message", slog.String("error", err.Error()))
		}
		return
	}

	if err := msg.Ack(false); err != nil {
		consumer.log.Error("failed to ack message", slog.String("error", err.Error()))
	} else {
		consumer.log.Debug("message successfully acked", slog.Uint64("delivery_tag", msg.DeliveryTag))
	}
}
