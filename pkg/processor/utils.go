package processor

import (
	"time"

	"github.com/occmundial/consumer-clevertap-applies/pkg/kafka"
)

// secondsToTimeDuration : convierte de segundos a time.Duration
func secondsToTimeDuration(seconds int) time.Duration {
	return time.Second * time.Duration(seconds)
}

func newProcessStatus[T any](status string, message kafka.MessageForRead[T], err error) ProcessStatus[T] {
	return ProcessStatus[T]{
		Status:  status,
		Message: message,
		Error:   err,
	}
}

func mapMessageToRetryEvent[T any](message T, topic string) MessageToSend[T] {
	return MessageToSend[T]{
		Topic:   topic,
		Message: message,
	}
}
