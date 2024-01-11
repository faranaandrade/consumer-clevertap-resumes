package kafka

import (
	"github.com/segmentio/kafka-go"

	"time"
)

// getKafkaReader :
func getKafkaReader(config *QueueSetup) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		GroupID:        config.GroupID,
		Topic:          config.TopicMain,
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
	})
}

func kafkaMessageToMessageForRead[T any](kafkaMessage *kafka.Message, location *time.Location,
	deserializer MessageDeserializer[T]) (MessageForRead[T], error) {
	deserializedMessage, err := deserializer.GetMessageFromBytes(kafkaMessage.Value)

	return MessageForRead[T]{
		Key:           string(kafkaMessage.Key),
		Topic:         kafkaMessage.Topic,
		Time:          kafkaMessage.Time.In(location),
		Message:       deserializedMessage,
		FlatMessage:   string(kafkaMessage.Value),
		SourceMessage: kafkaMessage,
	}, err
}
