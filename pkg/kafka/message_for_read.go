package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// MessageForRead :
type MessageForRead[T any] struct {
	Key           string         `json:"key" example:"c4bf18e8-10f6-48ca-aa6f-61c7e8c25578"`
	Topic         string         `json:"topic"`
	Time          time.Time      `json:"time"`
	Message       T              `json:"message"`
	FlatMessage   string         `json:"flatMessage"`
	SourceMessage *kafka.Message `json:"-"`
}
