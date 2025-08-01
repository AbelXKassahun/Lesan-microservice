package rabbitmq

import (
	"encoding/json"
	"log"

	"gamification-service/internal/app"

	"github.com/streadway/amqp"
)

func StartLessonCompletedConsumer(rabbitURL string, xpService *app.XPService) error {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"lesson_completed", // queue name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var event app.LessonCompletedEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Println("❌ Failed to parse event:", err)
				continue
			}
			log.Printf("📥 Received LessonCompleted event: %+v\n", event)

			// XP/streak/badge logic
			if event.Event == "LessonCompleted" {
				xpService.HandleLessonCompleted(event)
			}

		}
	}()

	log.Println("✅ Listening for LessonCompleted events on queue `lesson_completed`")
	return nil
}
