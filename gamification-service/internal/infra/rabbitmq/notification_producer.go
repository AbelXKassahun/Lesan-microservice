package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type NotificationMessage struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Message  string   `json:"message"`
	Channels []string `json:"channels"`
}

type NotificationPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewNotificationPublisher(amqpURL, queueName string) (*NotificationPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &NotificationPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

/*
  - msg := NotificationMessage{
    UserID:   "123",
    Email:    "michaelberhan08@gmail.com",
    Message:  "fsefsefse",
    Channels: []string{"email"},
    }
*/
func (p *NotificationPublisher) Publish(msg NotificationMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	err = p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("✅ Notification published to queue '%s' for user %s", p.queue.Name, msg.UserID)
	return nil
}

func (p *NotificationPublisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
