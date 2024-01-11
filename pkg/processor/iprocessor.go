package processor

// IProcessor :
type IProcessor[T any] interface {
	CheckTargetHealth() error
	ProcessMessage() ProcessStatus[T]
}
