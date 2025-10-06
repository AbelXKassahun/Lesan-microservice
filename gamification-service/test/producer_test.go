package test

// go test -count=1 ./internal/infra/rabbitmq -v
// or
// go clean -testcache
// go test ./internal/infra/rabbitmq -v
import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

type LessonCompletedEvent struct {
	Event     string `json:"event"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	LessonID  string `json:"lesson_id"`
	Timestamp string `json:"timestamp"`
	XP        int    `json:"xp"`
}

func TestPublishLessonCompletedEvent(t *testing.T) {
	godotenv.Load("../../../.env")
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Println("❌ RABBITMQ_URL is not set. Using default value: amqp://guest:guest@localhost:5672/")
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		t.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("❌ Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Make sure queue exists
	_, err = ch.QueueDeclare(
		"lesson_completed",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("❌ Failed to declare queue: %v", err)
	}

	event := LessonCompletedEvent{
		Event: "LessonCompleted",
		// UserID:    "5c8530b5-5056-466c-b82f-82874cf226c4",
		UserID:    "u124",
		Email:     "kayanaloi32@gmail.com",
		LessonID:  "lessonA",
		Timestamp: time.Now().Format(time.RFC3339),
		XP:        15,
	}

	body, _ := json.Marshal(event)

	err = ch.Publish(
		"",                 // exchange
		"lesson_completed", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		t.Fatalf("❌ Failed to publish message: %v", err)
	}

	log.Println("📤 Published LessonCompleted")
}
