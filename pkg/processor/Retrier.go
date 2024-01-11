package processor

// Retryer :
type Retryer interface {
	CheckTargetHealth() error
	CreateEvent(message any) error
}
