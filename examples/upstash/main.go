package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	normal := KafkaNormal{
		Broker:         os.Getenv("KAFKA_BROKER"),
		SchemaRegistry: os.Getenv("KAFKA_SCHEMA_REGISTRY"),
		Username:       os.Getenv("KAFKA_USERNAME"),
		Password:       os.Getenv("KAFKA_PASSWORD"),
	}

	go normal.Producer(ctx, os.Getenv("KAFKA_TOPIC"))
	go normal.Consumer(ctx, os.Getenv("KAFKA_TOPIC"), "algo")

	// Create a channel to receive signals
	sigChan := make(chan os.Signal, 1)

	// Register the signal handler
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	<-sigChan
	cancel()
}
