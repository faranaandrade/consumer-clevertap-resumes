package setup

import (
	"github.com/occmundial/consumer-clevertap-applies/config"
	"github.com/occmundial/consumer-clevertap-applies/pkg/events"
)

func NewEventsSetup(configuration *config.Configuration) *events.Setup {
	return &events.Setup{
		APIEvents:  configuration.APIEvents,
		APITimeout: configuration.APITimeout,
		HTTPClient: configuration.HTTPClient,
	}
}
