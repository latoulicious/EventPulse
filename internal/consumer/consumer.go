package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	c *kafka.Consumer
}

func NewConsumer(brokers []string, topic string) *Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	topics := []string{topic}
	c.SubscribeTopics(topics, nil)

	return &Consumer{c: c}
}

func (c *Consumer) Start() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	run := true

	for run {
		select {
		case sig := <-signals:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := c.c.ReadMessage(-1)
			if err == nil {
				fmt.Printf("Consumed message from topic %s: %s\n", msg.TopicPartition.Topic, string(msg.Value))
			} else {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.c.Close()
}
