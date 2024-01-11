package kafka

// QueueGetter :
type QueueGetter[T any] interface {
	GetMessage() (MessageForRead[T], error)
}
