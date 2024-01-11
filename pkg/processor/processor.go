package processor

import (
	"context"
	"time"

	"github.com/occmundial/consumer-clevertap-applies/pkg/kafka"
	"github.com/occmundial/consumer-clevertap-applies/pkg/location"

	"github.com/occmundial/go-common/logger"
)

// Processor :
type Processor[T any] struct {
	QueueGetter     kafka.QueueGetter[T]
	Retryer         Retryer
	MessageProcesor MessageResolver[T]
	Configuration   *ProcessingRetriesSetup
	Location        location.Locater
	Log             logger.Logger
}

// CheckTargetHealth :
func (processor *Processor[T]) CheckTargetHealth() error {
	return processor.Retryer.CheckTargetHealth()
}

// ProcessMessage :
func (processor *Processor[T]) ProcessMessage() ProcessStatus[T] {
	var messageToProcess, err = processor.QueueGetter.GetMessage()
	if err != nil {
		// error durante la lectura del mensaje: se administra como poison message fuera del consumidor
		return newProcessStatus[T](StatusReadMessageError, messageToProcess, err)
	}
	// lectura exitosa del mensaje de la cola
	if !processor.MessageProcesor.IsValid(&messageToProcess.Message) {
		return newProcessStatus[T](StatusInvalidMessage, messageToProcess, nil)
	}
	status, err := processor.processMessage(context.Background(), &messageToProcess)
	return newProcessStatus[T](status, messageToProcess, err)
}

func (processor *Processor[T]) processMessage(ctx context.Context, message *kafka.MessageForRead[T]) (string, error) {
	if wait := processor.calculateWaitingTime(message.Time, processor.Location.DateNow()); wait > 0 {
		processor.Log.Debugf("wait %f seconds to process message", wait.Seconds())
		time.Sleep(wait)
	}
	if err := processor.MessageProcesor.Process(ctx, &message.Message); err != nil {
		// Error en procesamiento del mensaje: mandar a tópico de reintentos
		processor.sendRetryMessage(&message.Message)
		return StatusProcessError, err
	}
	return StatusFullProcessOK, nil
}

// calculateWaitingTime: tiempo a esperar para procesar el mensaje
func (processor *Processor[T]) calculateWaitingTime(messageTime, currentTime time.Time) time.Duration {
	if processor.Configuration.WaitForRetry > 0 {
		processTime := messageTime.Add(secondsToTimeDuration(processor.Configuration.WaitForRetry))
		if timeToProcess := processTime.Sub(currentTime); timeToProcess > 0 {
			return timeToProcess
		}
	}
	return 0
}

func (processor *Processor[T]) sendRetryMessage(message *T) {
	if retryTopic := processor.getRetryTopic(message); len(retryTopic) > 0 {
		event := mapMessageToRetryEvent(*message, retryTopic)
		if err := processor.Retryer.CreateEvent(event); err != nil {
			processor.Log.Error("services", "sendRetryMessage", err)
		} else {
			processor.Log.Debugf("reprocessing message sent to the topic '%s'", retryTopic)
		}
	}
}

func (processor *Processor[T]) getRetryTopic(message *T) string {
	if processor.Configuration.MaxRetries > 0 {
		// cantidad máxima de reprocesamientos permitidos en el tópico principal es mayor que cero
		retryNumber := processor.MessageProcesor.GetRetryNumber(message)
		if retryNumber < processor.Configuration.MaxRetries {
			// aún quedan reprocesamientos disponibles en el tópico principal
			processor.MessageProcesor.SetRetryNumber(message, retryNumber+1)
			return processor.Configuration.TopicMain
		}
	}
	// no se permiten reintentos en el tópico principal O ya no hay reintentos disponibles
	if len(processor.Configuration.TopicRetry) > 0 {
		// inicializo la cantidad de reintentos consumidos en cero para el tópico de reintentos
		processor.MessageProcesor.SetRetryNumber(message, 0)
		return processor.Configuration.TopicRetry
	}
	return ""
}
