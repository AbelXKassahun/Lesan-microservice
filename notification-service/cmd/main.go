package main

import (
    "log"
    rabbitmq "notification-service/internal/infra/RabbitMQ"
    "notification-service/internal/infra/RabbitMQ/workers"
)

func main() {
    conn, ch := rabbitmq.ConnectRabbit("amqp://guest:guest@localhost:5672/")
    defer conn.Close()
    defer ch.Close()

    queueName := "notification.queue"

    _, err := ch.QueueDeclare(
        queueName,
        true,  // durable
        false, // auto-delete
        false, // exclusive
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    workers.StartConsumer(ch, queueName)
}
