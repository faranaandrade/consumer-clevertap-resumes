package setup

import (
	"github.com/occmundial/consumer-clevertap-applies/config"
	"github.com/occmundial/consumer-clevertap-applies/pkg/kafka"
)

func NewQueueSetup(configuration *config.Configuration) *kafka.QueueSetup {
	return &kafka.QueueSetup{
		TopicMain:    configuration.QueueSetup.TopicMain,
		GroupID:      configuration.QueueSetup.GroupID,
		Brokers:      configuration.QueueSetup.Brokers,
		Timeout:      configuration.QueueSetup.Timeout,
		RequestDelay: configuration.QueueSetup.RequestDelay,
	}
}
