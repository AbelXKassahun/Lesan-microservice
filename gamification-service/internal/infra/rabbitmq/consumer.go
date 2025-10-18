package rabbitmq

import (
	"encoding/json"
	"log"

	"gamification-service/internal/app"
	"gamification-service/internal/domain"

	"github.com/streadway/amqp"
)

type LessonCompletedConsumer struct {
	XPService     *app.XPService
	StreakService *app.StreakService
	BadgeService  *app.BadgeService
}

func NewLessonCompletedConsumer(xpService *app.XPService,
	streakService *app.StreakService,
	badgeService *app.BadgeService) *LessonCompletedConsumer {
	return &LessonCompletedConsumer{
		XPService:     xpService,
		StreakService: streakService,
		BadgeService:  badgeService,
	}
}

func (c *LessonCompletedConsumer) StartLessonCompletedConsumer(rabbitURL string) error {
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
			var event domain.LessonCompletedEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Println("❌ Failed to parse event:", err)
				continue
			}
			log.Printf("📥 Received LessonCompleted event: %+v\n", event)

			// XP/streak/badge logic
			if event.Event == "LessonCompleted" {
				// this whole block of code needs to be refactored, logic should be handled in another file, this file is only for consuming events

				// notify on lesson complete
				notify(event.UserID, event.Email, "Congrats on completing a lesson, keep up the good work")

				c.XPService.UpdateXP(event.UserID, event.XP)
				streak, _ := c.StreakService.UpdateStreak(event.UserID)
				// updateProfile()
				switch streak {
				case 1:
					badgeAwarded, err := c.BadgeService.TryAwardBadge(event.UserID, domain.FirstLesson.String())
					if err != nil {
						log.Println("Couldn't award badge: ", err)
					}

					c.XPService.UpdateXP(event.UserID, 100)

					if badgeAwarded {
						notify(event.UserID, event.Email, "Congrats on finishing your first lesson, you have been awarded the 'FirstLesson' badge. Keep at it.")
					}
				case 7:
					badgeAwarded, err := c.BadgeService.TryAwardBadge(event.UserID, domain.FirstWeek.String())
					if err != nil {
						log.Println("Couldn't award badge: ", err)
					}

					c.XPService.UpdateXP(event.UserID, 250)

					if badgeAwarded {
						notify(event.UserID, event.Email, "A week straight of learning amharic, congrats. You have been awarded the 'FirstWeek' badge. We hope to see more of you.")
					}
				}
				c.updateProfile(event.UserID)
			}
		}
	}()

	log.Println("✅ Listening for LessonCompleted events on queue `lesson_completed`")
	return nil
}

func notify(userID, email, message string) {
	p, err := NewNotificationPublisher("amqp://guest:guest@localhost:5672/", "notification.queue")
	if err != nil {
		log.Fatalf("Error creating publisher: %v", err)
	}
	defer p.Close()

	msg := NotificationMessage{
		UserID:   userID,
		Email:    email,
		Message:  message,
		Channels: []string{"email"},
	}

	if err := p.Publish(msg); err != nil {
		log.Fatalf("Publish error: %v", err)
	}
	log.Println("✔ Notification Sent: ", message)
}

func (c *LessonCompletedConsumer) updateProfile(userID string) {
	p, err := NewUpdateProfileProducer("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error creating publisher: %v", err)
	}
	defer p.Close()

	userXP, err := c.XPService.GetXPByUserID(userID)
	if err != nil {
		log.Fatalf("Publish error: %v", err)
	}

	streak, err := c.StreakService.GetStreakByUserID(userID)
	if err != nil {
		log.Fatalf("Publish error: %v", err)
	}

	msg := UpdateProfileEvent{
		UserID:           userID,
		Streak:           streak.CurrentStreak,
		XP:               userXP.Total,
		LessonsCompleted: streak.ActiveDays,
		League:           string(userXP.League),
	}

	if err := p.Publish(msg); err != nil {
		log.Fatalf("Publish error: %v", err)
	}
	log.Println("✔ Update profile event sent")
}
