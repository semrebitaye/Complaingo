package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection // connection to rabbitmq
	channel *amqp.Channel    //lightweight communication channel on top of connection
	queue   amqp.Queue       //queue's publishing to
}

// newProducer connect to rabbitmq and sets up the queue
func NewProducer(amqpURL, queueName string) *Producer {
	// connect to rabbitmq
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %v", err)
	}

	// open channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// create queue if not exist
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}

	return &Producer{conn: conn, channel: ch, queue: q}
}

// publish a message to the queue
func (p *Producer) SendMessage(body string) error {
	err := p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message sent to queue: %s", body)
	return nil
}
