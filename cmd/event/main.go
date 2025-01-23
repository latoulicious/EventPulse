package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"eventpulse/internal/config"
	"eventpulse/internal/consumer"
	"eventpulse/internal/producer"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize producer and consumer
	prod := producer.NewProducer(cfg.KafkaBrokers)
	cons := consumer.NewConsumer(cfg.KafkaBrokers, cfg.Topic)

	// Start producer and consumer
	go prod.Start(cfg.Topic)
	go cons.Start()

	// Wait for termination signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	log.Println("Shutting down...")

	prod.Close()
	cons.Close()
	log.Println("Shutdown complete")
}
