package rabbitmq

import (
    "log"
    "github.com/streadway/amqp"
)

func ConnectRabbit(url string) (*amqp.Connection, *amqp.Channel) {
    conn, err := amqp.Dial(url)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open channel: %v", err)
    }

    return conn, ch
}