package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)
type KafkaConsumer struct {
	Brokers []string
	Topic   string
	GroupID string
}
func NewKafkaConsumer(brokers []string, topic, gropID string) *KafkaConsumer {
	return &KafkaConsumer{
		Brokers: brokers,
		Topic:   topic,
		GroupID: gropID,
	}
}
func (kc *KafkaConsumer) StartConsuming(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(kc.Brokers, kc.GroupID, config)
	if err != nil {
		log.Fatal("Failed to create consumer group:", err)
	}
	handler := consumerGroupHandler{}
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{kc.Topic}, handler); err != nil {
				log.Println("Kafka consume error:", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
}
type consumerGroupHandler struct{}
func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("kafka message received: %s\n", string(message.Value))
		session.MarkMessage(message, "")
	}
	return nil
}
