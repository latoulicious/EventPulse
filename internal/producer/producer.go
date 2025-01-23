package producer

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	p *kafka.Producer
}

func NewProducer(broker string) *Producer { // Change to accept a single string
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	return &Producer{p: p}
}

func (p *Producer) Start(topic string) {
	defer p.p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Simulate event generation
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Event %d", i)
		p.p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, nil)
		time.Sleep(1 * time.Second) // Simulate delay between events
	}

	// Wait for message deliveries before shutting down
	p.p.Flush(15 * 1000)
}

func (p *Producer) Close() {
	p.p.Close()
}
