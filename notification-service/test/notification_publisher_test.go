package test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/streadway/amqp"
	"notification-service/internal/domain"
)

func TestPublishNotification(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	queueName := "notification.queue"

	notification := domain.Notification{
		UserID:   "123",
		Email:    "michaelberhan08@gmail.com",
		Message:  "Ato Sefineh, ante wesha, unsubscripe from this app",
		Channels: []string{"email"},
	}

	body, err := json.Marshal(notification)
	if err != nil {
		t.Fatalf("Failed to marshal notification: %v", err)
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	log.Println("✅ Notification published from test")
}
