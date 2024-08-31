package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KafkaNormal struct {
	Broker         string
	SchemaRegistry string
	Username       string
	Password       string
}

func (k KafkaNormal) Producer(ctx context.Context, topic string) error {
	mechanism, _ := scram.Mechanism(scram.SHA256, k.Username, k.Password)
	w := kafka.Writer{
		Addr:  kafka.TCP(k.Broker),
		Topic: topic,
		Transport: &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		},
	}
	defer w.Close()

	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			value := NewUser()
			data, err := json.Marshal(value)
			if err != nil {
				panic(err)
			}
			err = w.WriteMessages(ctx, kafka.Message{Value: data})
			if err != nil {
				panic(err)
			}
			slog.Info("producer", "value", string(data))
		}
		time.Sleep(time.Second * 3)
	}

	return nil
}

func (k KafkaNormal) Consumer(ctx context.Context, topic string, groupId string) {
	mechanism, _ := scram.Mechanism(scram.SHA512, k.Username, k.Password)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.Broker},
		GroupID: groupId,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
	})
	defer r.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*120) // Increase the timeout
	defer cancel()

	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			message, _ := r.ReadMessage(ctx)
			slog.Info("consumer", "partition", message.Partition, "offset", message.Offset, "value", string(message.Value))
		}
	}
}
