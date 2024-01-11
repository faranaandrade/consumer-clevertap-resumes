package setup

import (
	"github.com/occmundial/consumer-clevertap-applies/internal/app/updater"
	"github.com/occmundial/consumer-clevertap-applies/internal/models"
	"github.com/occmundial/consumer-clevertap-applies/pkg/events"
	"github.com/occmundial/consumer-clevertap-applies/pkg/kafka"
	"github.com/occmundial/consumer-clevertap-applies/pkg/location"
	"github.com/occmundial/consumer-clevertap-applies/pkg/processor"
	"github.com/occmundial/go-common/logger"
)

// NewConsumerService : Factory que crea un "Processor"
func NewConsumerService(
	queueGetter *kafka.QueueGetterKafka[models.MessageToProcess],
	retrier *events.Repository[models.MessageToProcess],
	messageProcesor *updater.MessageResolverUpdater, loc *location.Location,
	configuration *processor.ProcessingRetriesSetup, log *logger.Log) *processor.Processor[models.MessageToProcess] {
	return &processor.Processor[models.MessageToProcess]{
		QueueGetter:     queueGetter,
		Retryer:         retrier,
		Configuration:   configuration,
		MessageProcesor: messageProcesor,
		Location:        loc,
		Log:             log,
	}
}
