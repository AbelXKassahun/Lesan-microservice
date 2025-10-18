package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// UpdateProfileEvent is the payload structure expected by the update consumer
type UpdateProfileEvent struct {
	Event            string `json:"event"`
	UserID           string `json:"user_id"`
	Streak           int    `json:"streak,omitempty"`
	XP               int    `json:"xp,omitempty"`
	LessonsCompleted int    `json:"lessons_completed,omitempty"`
	League           string `json:"league,omitempty"`
}

type UpdateProfileProducer struct {
	channel *amqp.Channel
	queue   amqp.Queue
	conn    *amqp.Connection
}

func NewUpdateProfileProducer(rabbitURL string) (*UpdateProfileProducer, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"update_profile", // same queue as consumer
		true,             // durable
		false,            // auto-delete
		false,            // exclusive
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &UpdateProfileProducer{
		channel: ch,
		queue:   q,
		conn:    conn,
	}, nil
}

func (p *UpdateProfileProducer) Publish(event UpdateProfileEvent) error {
	event.Event = "update_profile"

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",           // default exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("📤 Sent update_profile event: %+v\n", event)
	return nil
}

func (p *UpdateProfileProducer) Close() {
	p.channel.Close()
	p.conn.Close()
}
