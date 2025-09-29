package rabbitmq

// import (
// 	"encoding/json"
// 	"log"

// 	"lesson-service/internal/app"
// 	"lesson-service/internal/domain"

// 	"github.com/streadway/amqp"
// )

// type LessonCompletedConsumer struct {
// 	XPService     *app.XPService
// 	StreakService *app.StreakService
// 	BadgeService  *app.BadgeService
// }


// func NewLessonCompletedConsumer(xpService *app.XPService,
// 	streakService *app.StreakService,
// 	badgeService *app.BadgeService) *LessonCompletedConsumer {
// 	return &LessonCompletedConsumer{
// 		XPService:     xpService,
// 		StreakService: streakService,
// 		BadgeService:  badgeService,
// 	}
// }

// func (c *LessonCompletedConsumer) StartLessonCompletedConsumer(rabbitURL string) error {
// 	conn, err := amqp.Dial(rabbitURL)
// 	if err != nil {
// 		return err
// 	}
// 	ch, err := conn.Channel()
// 	if err != nil {
// 		return err
// 	}

// 	q, err := ch.QueueDeclare(
// 		"lesson_completed", // queue name
// 		true,               // durable
// 		false,              // delete when unused
// 		false,              // exclusive
// 		false,              // no-wait
// 		nil,                // args
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	msgs, err := ch.Consume(
// 		q.Name,
// 		"",
// 		true,  // auto-ack
// 		false, // exclusive
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	go func() {
// 		for d := range msgs {
// 			var event domain.LessonCompletedEvent
// 			err := json.Unmarshal(d.Body, &event)
// 			if err != nil {
// 				log.Println("❌ Failed to parse event:", err)
// 				continue
// 			}
// 			log.Printf("📥 Received LessonCompleted event: %+v\n", event)

// 			// XP/streak/badge logic
// 			if event.Event == "LessonCompleted" {
// 				// this (the block of code in this `if statement`) needs to be refactored, logic should be handled i n another file, this file is only for consuming events
// 				c.XPService.UpdateXP(event.UserID, event.XP)
// 				streak, _ := c.StreakService.UpdateStreak(event.UserID)
// 				switch streak {
// 				case 1:
// 					c.BadgeService.TryAwardBadge(event.UserID, domain.FirstLesson.String())
// 					c.XPService.UpdateXP(event.UserID, 100)
// 				case 7:
// 					c.BadgeService.TryAwardBadge(event.UserID, domain.FirstWeek.String())
// 					c.XPService.UpdateXP(event.UserID, 250)
// 				}
// 			}

// 		}
// 	}()

// 	log.Println("✅ Listening for LessonCompleted events on queue `lesson_completed`")
// 	return nil
// }
