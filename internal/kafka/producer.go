package kafka

import (
	"log"
	"strings"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_5_0_0 //set kafka version

	// create admin client to check/create topic
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatal("Failed to create kafka admin:", err)
	}
	defer admin.Close()

	// try to create the topic if not exists
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil && !strings.Contains(err.Error(), "Topic with this name already exists") {
		log.Fatal("Failed to create topic:", err)
	} else {
		log.Println("kafka topic ready", topic)
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal("Failed to create kafka producer:", err)
	}

	return &KafkaProducer{
		Producer: producer,
		Topic:    topic,
	}
}

func (kp *KafkaProducer) SendMessage(message string) {
	msg := &sarama.ProducerMessage{
		Topic: kp.Topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := kp.Producer.SendMessage(msg)
	if err != nil {
		log.Println("kafka send failed:", err)
	} else {
		log.Println("kafka message sent:", message)
	}
}
