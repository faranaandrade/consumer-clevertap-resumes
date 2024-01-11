package processor

import "github.com/occmundial/consumer-clevertap-applies/pkg/kafka"

// ProcessStatus :
type ProcessStatus[T any] struct {
	Status  string `json:"status" example:"full-processed-message"` // utils.StatusFullProcessOK
	Message kafka.MessageForRead[T]
	Error   error
}
