package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type KafkaSchema struct {
	Broker         string
	SchemaRegistry string
	Username       string
	Password       string
}

func (k KafkaSchema) Producer(ctx context.Context, topic string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": k.Broker,
		"sasl.mechanism":    "SCRAM-SHA-256",
		"security.protocol": "SASL_SSL",
		"sasl.username":     k.Username,
		"sasl.password":     k.Password,
	})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		k.SchemaRegistry,
		k.Username,
		k.Password,
	))
	if err != nil {
		panic(err)
	}

	ser, err := avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		panic(err)
	}

	deliveryChan := make(chan kafka.Event)
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					slog.Info("schema producer delivery failed",
						"TopicPartition", ev.TopicPartition,
					)
				} else {
					slog.Info("schema producer delivered message",
						"TopicPartition", ev.TopicPartition,
					)
				}
			}
		}
	}()

	newMsg := func() error {
		value := NewUser()
		payload, err := ser.Serialize(topic, &value)
		if err != nil {
			panic(err)
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte("test"),
			Value:          payload,
		}, deliveryChan)
		if err != nil {
			panic(err)
		}
		return nil
	}

	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			if err := newMsg(); err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second * 3)
	}

	return nil
}

func (k KafkaSchema) Consumer(ctx context.Context, topic string, groupId string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": k.Broker,
		"sasl.mechanism":    "SCRAM-SHA-256",
		"security.protocol": "SASL_SSL",
		"sasl.username":     k.Username,
		"sasl.password":     k.Password,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		k.SchemaRegistry,
		k.Username,
		k.Password,
	))
	if err != nil {
		panic(err)
	}

	deser, err := avro.NewGenericDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())
	if err != nil {
		panic(err)
	}

	err = c.Subscribe(topic, nil)
	if err != nil {
		panic(err)
	}

	procMsg := func() {
		ev := c.Poll(5000)
		if ev == nil {
			slog.Info("schema consumer", "topic", topic, "ev", ev)
			return
		}

		switch e := ev.(type) {
		case *kafka.Message:
			value := User{}
			deser.DeserializeInto(*e.TopicPartition.Topic, e.Value, &value)
			slog.Info("schema consumer", "TopicPartition", ev.String(), "key", e.Key, "val", value)
		case kafka.Error:
			panic(e)
		default:
			slog.Info("Ignored", "ev", e)
		}
	}

	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			procMsg()
		}
	}
}
