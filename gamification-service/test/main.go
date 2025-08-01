package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type LessonCompletedEvent struct {
	Event     string `json:"event"`
	UserID    string `json:"user_id"`
	LessonID  string `json:"lesson_id"`
	Timestamp string `json:"timestamp"`
	XP        int    `json:"xp"`
}

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	defer conn.Close()
	defer ch.Close()

	event := LessonCompletedEvent{
		Event:     "LessonCompleted",
		UserID:    "u123",
		LessonID:  "lessonA",
		Timestamp: time.Now().Format(time.RFC3339),
		XP:        15,
	}

	body, _ := json.Marshal(event)

	ch.Publish(
		"",                 // exchange
		"lesson_completed", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	log.Println("📤 Published LessonCompleted")
}
