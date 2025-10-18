package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"profile-service/internal/app"
	"profile-service/internal/domain"
	"time"

	// "time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type CreateProfileConsumer struct {
	Service app.ProfileService
}

func NewCreateProfileConsumer(service app.ProfileService) *CreateProfileConsumer {
	return &CreateProfileConsumer{Service: service}
}

func (c *CreateProfileConsumer) StartCreateProfileConsumer(rabbitURL string) error {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"create_profile", // queue name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // args
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
			var event struct {
				Event  string `json:"event"`
				UserID string `json:"user_id"`
				Email  string `json:"email"`
			}

			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Println("❌ Failed to parse event:", err)
				continue
			}

			if event.Event != "create_profile" {
				continue
			}

			log.Printf("📥 Received CreateProfile event: %+v\n", event)

			userID, err := uuid.Parse(event.UserID)
			if err != nil {
				log.Println("❌ Invalid user_id:", err)
				continue
			}

			profile := &domain.Profile{
				UserID:      userID,
				Email:       event.Email,
				MemberSince: time.Now(),
			}

			// ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
			ctx := context.Background()

			if err := c.Service.CreateProfile(ctx, profile); err != nil {
				log.Println("❌ Failed to create profile:", err)
				continue
			}

			log.Println("✅ Profile created successfully for user:", event.UserID)
		}
	}()

	log.Println("✅ Listening for CreateProfile events on queue `create_profile`")
	return nil
}
