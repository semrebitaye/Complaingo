package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection // connection to rabbitmq
	channel *amqp.Channel    //lightweight communication channel on top of connection
	queue   amqp.Queue       //queue's publishing to
}

// function to craete a new rabbitmq connection and channel
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
	q, err := ch.QueueDeclare(
		queueName, //name
		true,      //durable
		false,     //delete when unused
		false,     //exclusive
		false,     //no wait
		nil,       //arguments
	)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}

	return &Producer{conn: conn, channel: ch, queue: q}
}

// publish a message to the queue
func (p *Producer) SendMessage(body string) error {
	// convert the message to json format
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to marshal the message: %v", err)
		return err
	}

	// publish the meassage to the queue
	err = p.channel.Publish(
		"",           //excahnge
		p.queue.Name, //routing key(queue name)
		false,        //mandatory
		false,        //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonBody,
		})
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message sent to queue: %s", body)
	return nil
}
