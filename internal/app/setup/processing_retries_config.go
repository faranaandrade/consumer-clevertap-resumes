package setup

import (
	"github.com/occmundial/consumer-clevertap-resumes/config"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/processor"
)

func NewProcessingRetriesSetup(configuration *config.Configuration) *processor.ProcessingRetriesSetup {
	return &processor.ProcessingRetriesSetup{
		TopicRetry:   configuration.ProcessingRetries.TopicRetry,
		WaitForRetry: configuration.ProcessingRetries.WaitForRetry,
		MaxRetries:   configuration.ProcessingRetries.MaxRetries,
		TopicMain:    configuration.QueueSetup.TopicMain,
	}
}
