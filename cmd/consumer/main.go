package main

import (
	"github.com/occmundial/consumer-clevertap-resumes/cmd/di"
	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/consumer"
	"github.com/occmundial/go-common/logger"
)

func main() {
	getWorker().Run()
}

func getWorker() *consumer.Consumer[models.MessageToProcess] {
	var worker *consumer.Consumer[models.MessageToProcess]

	container := di.GetContainer()
	err := container.Invoke(func(p *consumer.Consumer[models.MessageToProcess]) {
		worker = p
	})

	if err != nil {
		logger.GetLogger().Fatal("main", "getWorker", err)
	}
	return worker
}
