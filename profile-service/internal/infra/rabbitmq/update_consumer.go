package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"profile-service/internal/domain"
	"profile-service/internal/app"
	// "time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type UpdateProfileConsumer struct {
	Service app.ProfileService
}

func NewUpdateProfileConsumer(service app.ProfileService) *UpdateProfileConsumer {
	return &UpdateProfileConsumer{Service: service}
}

func (c *UpdateProfileConsumer) StartUpdateProfileConsumer(rabbitURL string) error {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"update_profile", // queue name
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
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
				Event           string `json:"event"`
				UserID          string `json:"user_id"`
				Streak          int    `json:"streak"`
				XP              int    `json:"xp"`
				LessonsCompleted int   `json:"lessons_completed"`
				League          string `json:"league"`
			}

			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Println("❌ Failed to parse event:", err)
				continue
			}

			if event.Event != "update_profile" {
				continue
			}

			log.Printf("📥 Received UpdateProfile event: %+v\n", event)

			userID, err := uuid.Parse(event.UserID)
			if err != nil {
				log.Println("❌ Invalid user_id:", err)
				continue
			}

			profile := &domain.Profile{
				UserID:           userID,
				Streak:           event.Streak,
				XP:               event.XP,
				LessonsCompleted: event.LessonsCompleted,
				League:           event.League,
			}

			// /ctx, _ := context.WithTimeout(context.Background(), 2*time.Minute)
			ctx := context.Background()

			if err := c.Service.UpdateProfile(ctx, profile, userID); err != nil {
				log.Println("❌ Failed to update profile:", err)
				continue
			}

			log.Println("✅ Profile updated successfully for user:", event.UserID)
		}
	}()

	log.Println("✅ Listening for UpdateProfile events on queue `update_profile`")
	return nil
}
