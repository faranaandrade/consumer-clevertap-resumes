package kafka

import (
	"context"
	"time"

	"github.com/occmundial/consumer-clevertap-applies/pkg/location"
	"github.com/segmentio/kafka-go"
)

// NewQueueGetterKafka : Factory que crea un "QueueGetterKafka"
func NewQueueGetterKafka[T any](configuration *QueueSetup, l *location.Location,
	deserializer MessageDeserializer[T]) *QueueGetterKafka[T] {
	cr := QueueGetterKafka[T]{
		Configuration:       configuration,
		Location:            l.Location,
		MessageDeserializer: deserializer,
		KafkaTimeout:        time.Duration(configuration.Timeout) * time.Second,
		Reader:              getKafkaReader(configuration),
	}
	return &cr
}

// QueueGetterKafka :
type QueueGetterKafka[T any] struct {
	Configuration       *QueueSetup
	Location            *time.Location
	MessageDeserializer MessageDeserializer[T]
	Reader              *kafka.Reader
	KafkaTimeout        time.Duration
}

// GetMessage :
func (r QueueGetterKafka[T]) GetMessage() (MessageForRead[T], error) {
	message := MessageForRead[T]{}
	ctx, cancel := context.WithTimeout(context.Background(), r.KafkaTimeout)
	defer cancel()
	kafkaMessage, err := r.Reader.ReadMessage(ctx)
	if err == nil {
		message, err = kafkaMessageToMessageForRead[T](&kafkaMessage, r.Location, r.MessageDeserializer)
	}
	return message, err
}
