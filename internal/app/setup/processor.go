package setup

import (
	"github.com/occmundial/consumer-clevertap-resumes/pkg/consumer"
	"github.com/occmundial/go-common/logger"

	"github.com/occmundial/consumer-clevertap-resumes/config"
	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/processor"
)

// NewProcessor : Factory que crea un "Consumer"
func NewProcessor(configuration *config.Configuration, service *processor.Processor[models.MessageToProcess],
	log *logger.Log) *consumer.Consumer[models.MessageToProcess] {
	return &consumer.Consumer[models.MessageToProcess]{
		Configuration: &consumer.Setup{
			ArtifactVersion: configuration.ArtifactVersion,
			RequestDelay:    configuration.QueueSetup.RequestDelay,
		},
		Processor: service,
		Log:       log,
	}
}
