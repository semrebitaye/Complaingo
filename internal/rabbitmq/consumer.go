package rabbitmq

import (
	"crud_api/internal/domain/models"
	"crud_api/internal/websocket"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn   *amqp.Connection
	cahhel *amqp.Channel
	queue  amqp.Queue
}

func NewConsumer(amqpUrl, queueName string) *Consumer {
	// connect to rabbitmq
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %v", err)
	}

	// open channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	// create queue
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	log.Println("Consumer connected to rabbitmq and queue declared")

	return &Consumer{conn: conn, cahhel: ch, queue: q}
}

func (c *Consumer) StartConsuming() {
	msgs, err := c.cahhel.Consume(c.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("Waiting for messages...")

	// process messages in a goroutine
	go func() {
		for msg := range msgs {
			log.Printf("Message recieved: %s", string(msg.Body))

			var notif models.NotificationMessage
			if err := json.Unmarshal(msg.Body, &notif); err == nil && notif.Type == "complaint_created" {
				log.Printf("Notify admins: user %d created a complaint %s", (notif.UserID), notif.Complient)

				go websocket.SendToAdmins(notif)
			}
		}
	}()
}
