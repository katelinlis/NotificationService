package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	brokers []string
}

func NewConsumer(brokers []string) *Consumer {
	return &Consumer{
		brokers: brokers,
	}
}

func (c *Consumer) Subscribe(context context.Context, topic, groupID string, handler func([]byte) error) {

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		GroupID: groupID,
		Topic:   topic,
		//StartOffset: kafka.FirstOffset,
	})

	println("Kafka consumer started for topic:", topic)

	go func() {
		for {
			println("Waiting for messages...")
			m, err := c.reader.ReadMessage(context)

			if err != nil {
				log.Printf("Kafka read error: %v", err)
				continue
			}

			if err := handler(m.Value); err != nil {
				log.Printf("Handler error: %v", err)
			}
		}
	}()
}
