package rabbitmq

import (
	"fmt"
	"log/slog"
	"media-processing-platform/server/internal/config"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionManager struct {
	cfg config.RabbitMQConfig
	log *slog.Logger

	mu   sync.RWMutex
	conn *amqp.Connection

	done chan struct{}
}

func NewConnectionManager(cfg config.RabbitMQConfig, log *slog.Logger) (*ConnectionManager, error) {
	connectionManager := &ConnectionManager{
		cfg:  cfg,
		log:  log,
		done: make(chan struct{}),
	}

	if err := connectionManager.connectWithRetry(); err != nil {
		return nil, err
	}

	go connectionManager.reconnectLoop()

	return connectionManager, nil
}

func (connectionManager *ConnectionManager) connectWithRetry() error {
	var err error
	backoff := connectionManager.cfg.Connection.InitialBackoff
	maxRetries := connectionManager.cfg.Connection.MaxRetries

	for i := 1; i <= maxRetries; i++ {
		conn, err := amqp.Dial(connectionManager.cfg.DSN())
		if err == nil {
			connectionManager.mu.Lock()
			connectionManager.conn = conn
			connectionManager.mu.Unlock()
			connectionManager.log.Info("successfully connected to RabbitMQ")
			return nil
		}

		connectionManager.log.Warn("failed to connect to RabbitMQ, retrying...",
			slog.Int("attempt", i),
			slog.String("error", err.Error()),
			slog.Duration("backoff", backoff),
		)

		time.Sleep(backoff)
		backoff *= 2
	}

	return fmt.Errorf("could not connect to RabbitMQ after %d retries: %w", maxRetries, err)
}

func (connectionManager *ConnectionManager) reconnectLoop() {
	for {
		connectionManager.mu.RLock()
		conn := connectionManager.conn
		connectionManager.mu.RUnlock()

		if conn == nil {
			return
		}

		closeErr := <-conn.NotifyClose(make(chan *amqp.Error))
		if closeErr == nil {
			return
		}

		connectionManager.log.Error("RabbitMQ connection lost! Reconnecting...",
			slog.String("reason", closeErr.Error()))

		for {
			select {
			case <-connectionManager.done:
				return
			default:
			}

			if err := connectionManager.connectWithRetry(); err == nil {
				connectionManager.log.Info("successfully reconnected to RabbitMQ")
				break
			}

			select {
			case <-connectionManager.done:
				return
			case <-time.After(connectionManager.cfg.Connection.ReconnectInterval):
			}
		}
	}
}

func (connectionManager *ConnectionManager) GetConnection() (*amqp.Connection, error) {
	connectionManager.mu.RLock()
	defer connectionManager.mu.RUnlock()

	if connectionManager.conn == nil || connectionManager.conn.IsClosed() {
		return nil, fmt.Errorf("rabbitmq connection is closed")
	}

	return connectionManager.conn, nil
}

func (connectionManager *ConnectionManager) Close() error {
	close(connectionManager.done)
	connectionManager.mu.Lock()
	defer connectionManager.mu.Unlock()

	if connectionManager.conn != nil && !connectionManager.conn.IsClosed() {
		return connectionManager.conn.Close()
	}
	return nil
}
