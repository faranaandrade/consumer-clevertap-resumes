package setup

import (
	"github.com/occmundial/consumer-clevertap-resumes/internal/app/updater"
	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/kafka"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/location"
)

func GetQueueGetterKafka(configuration *kafka.QueueSetup, locater *location.Location,
	deserializer *updater.MessageDeserializerMessageToProcess) *kafka.QueueGetterKafka[models.MessageToProcess] {
	return kafka.NewQueueGetterKafka[models.MessageToProcess](configuration, locater, deserializer)
}
