package setup

import (
	"github.com/occmundial/consumer-clevertap-applies/internal/app/updater"
	"github.com/occmundial/consumer-clevertap-applies/internal/models"
	"github.com/occmundial/consumer-clevertap-applies/pkg/kafka"
	"github.com/occmundial/consumer-clevertap-applies/pkg/location"
)

func GetQueueGetterKafka(configuration *kafka.QueueSetup, locater *location.Location,
	deserializer *updater.MessageDeserializerMessageToProcess) *kafka.QueueGetterKafka[models.MessageToProcess] {
	return kafka.NewQueueGetterKafka[models.MessageToProcess](configuration, locater, deserializer)
}
