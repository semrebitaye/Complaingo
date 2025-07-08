package rabbitmq

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/websocket"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
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
	q, err := ch.QueueDeclare(
		queueName, //queue name
		true,      //durable
		false,     //autodelete
		false,     //exclusive
		false,     //noWait
		nil,       //arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	log.Println("Consumer connected to rabbitmq and queue declared")

	return &Consumer{conn: conn, channel: ch, queue: q}
}

func (c *Consumer) StartConsuming() {
	// consume rabbitmq queue
	msgs, err := c.channel.Consume(
		c.queue.Name, //queue name
		"",           //consumer tag
		true,         //auto ack
		false,        //exclusive
		false,        //no-local
		false,        //no-wait
		nil,          //args
	)
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
