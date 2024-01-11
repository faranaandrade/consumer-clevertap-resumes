package processor

import "context"

type MessageResolver[T any] interface {
	Process(ctx context.Context, message *T) error
	IsValid(message *T) bool
	GetRetryNumber(message *T) int
	SetRetryNumber(message *T, value int)
}
