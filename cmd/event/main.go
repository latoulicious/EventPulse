package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/latoulicious/EventPulse/internal/config"
	"github.com/latoulicious/EventPulse/internal/consumer"
	"github.com/latoulicious/EventPulse/internal/producer"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize producer and consumer
	prod := producer.NewProducer(cfg.KafkaBrokers) // Pass the single string
	cons := consumer.NewConsumer([]string{cfg.KafkaBrokers}, cfg.Topic)

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
