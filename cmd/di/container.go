package di

import (
	"github.com/occmundial/consumer-clevertap-resumes/internal/app/updater"
	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/gorm"

	"sync"

	"github.com/occmundial/consumer-clevertap-resumes/pkg/events"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/location"
	"github.com/occmundial/go-common/logger"

	"github.com/occmundial/consumer-clevertap-resumes/config"
	"github.com/occmundial/consumer-clevertap-resumes/internal/app/setup"
	"go.uber.org/dig"
)

// https://blog.drewolson.org/dependency-injection-in-go

var (
	container *dig.Container
	once      sync.Once
)

// GetContainer :
func GetContainer() *dig.Container {
	once.Do(func() {
		container = buildContainer()
	})
	return container
}

// buildContainer :
func buildContainer() *dig.Container {
	c := dig.New()
	logger.GetLogger().FatalIfError("container", "buildContainer",
		c.Provide(logger.NewLog),
		c.Provide(location.NewLocation),
		c.Provide(config.NewConfiguration),
		c.Provide(setup.NewQueueSetup),
		c.Provide(setup.NewEventsSetup),
		c.Provide(setup.NewGormSetup),
		c.Provide(setup.NewProcessingRetriesSetup),
		c.Provide(updater.NewMessageDeserializer),
		c.Provide(events.NewEventsRepository[models.MessageToProcess]),
		c.Provide(setup.GetQueueGetterKafka),
		c.Provide(setup.NewConsumerService),
		c.Provide(setup.NewProcessor),
		c.Provide(gorm.NewDBGorm),
		c.Provide(updater.NewDataUserRepositoryDatabase),
		c.Provide(updater.NewMessageResolverUpdater),
		c.Provide(updater.NewAPIClevetapRepository))
	return c
}
