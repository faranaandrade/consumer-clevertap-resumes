package kafka

// MessageDeserializer :
type MessageDeserializer[T any] interface {
	GetMessageFromBytes(bytes []byte) (T, error)
}
