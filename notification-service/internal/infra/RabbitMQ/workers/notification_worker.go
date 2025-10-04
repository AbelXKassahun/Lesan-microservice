package workers

import (
    "encoding/json"
    "log"

    "github.com/streadway/amqp"
    "notification-service/internal/domain"
    "notification-service/internal/app"
)

func StartConsumer(ch *amqp.Channel, queueName string) {
    msgs, err := ch.Consume(
        queueName,
        "",
        true,  // auto-ack
        false, // not exclusive
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to consume messages: %v", err)
    }

    log.Printf("Listening for notifications on queue: %s", queueName)

    for msg := range msgs {
        var n domain.Notification
        if err := json.Unmarshal(msg.Body, &n); err != nil {
            log.Printf("Failed to unmarshal: %v", err)
            continue
        }

        log.Printf("Received notification for user %s", n.UserID)

        for _, ch := range n.Channels {
            switch ch {
            case "email":
                app.SendEmail(n.Email, n.Message)
            case "push":
                log.Println("[TODO] Push notification support not implemented yet")
            default:
                log.Printf("Unknown channel: %s", ch)
            }
        }
    }
}
