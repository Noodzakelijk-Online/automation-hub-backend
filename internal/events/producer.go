package events

import (
	"automation-hub-backend/internal/config"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type Publisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewPublisher(brokers []string, topic string) (*Publisher, error) {
	newConfig := sarama.NewConfig()
	newConfig.Producer.RequiredAcks = sarama.WaitForAll
	newConfig.Producer.Retry.Max = 5
	newConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, newConfig)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		producer: producer,
		topic:    topic,
	}, nil
}

func DefaultPublisher() *Publisher {
	producer, err := NewPublisher(config.AppConfig.Brokers, config.AppConfig.Topic)
	if err != nil {
		log.Fatalf("Failed to create default producer: %v", err)
	}
	return producer
}

func (p *Publisher) Close() error {
	return p.producer.Close()
}

func (p *Publisher) Publish(event *AutomationEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(message),
		Key:   sarama.StringEncoder(event.Automation.ID.String()),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Sent message to Kafka topic %s", p.topic)
	return nil
}
