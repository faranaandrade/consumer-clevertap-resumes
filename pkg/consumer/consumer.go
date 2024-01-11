package consumer

import (
	"context"
	"errors"
	"fmt"

	"github.com/occmundial/consumer-clevertap-resumes/pkg/processor"

	"time"

	"github.com/occmundial/consumer-clevertap-resumes/pkg/kafka"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/occmundial/go-common/logger"
)

// Consumer :
type Consumer[T any] struct {
	Configuration *Setup
	Processor     processor.IProcessor[T]
	Log           logger.Logger
}

// Run :
func (c *Consumer[T]) Run() {
	for {
		c.consume()
	}
}

func (c *Consumer[T]) consume() {
	fmt.Println()
	c.Log.Debugf(fmt.Sprintf("process started (version: %s)", c.Configuration.ArtifactVersion))
	if c.isTargetHealth() {
		for {
			if !c.processMessage() || !c.isTargetHealth() {
				break
			}
		}
	} else {
		c.Log.Debugf(StatusCheckProcessError)
	}
	c.Log.Debugf(fmt.Sprintf("process ended (version: %s)", c.Configuration.ArtifactVersion))
	c.wait()
}

func (c *Consumer[T]) isTargetHealth() bool {
	err := c.Processor.CheckTargetHealth()
	if err != nil {
		c.Log.Error("c", "isTargetHealth", err)
	}
	return err == nil
}

func (c *Consumer[T]) processMessage() bool {
	status := c.Processor.ProcessMessage()
	printStatus[T](status)
	return isHealthyStatus[T](status)
}

// isHealthyStatus : sin error y la cola tiene al menos un mensaje y kafka está saludable
func isHealthyStatus[T any](status processor.ProcessStatus[T]) bool {
	return status.Error == nil && status.Status != StatusProcessStartError
}

func printStatus[T any](status processor.ProcessStatus[T]) {
	isError, value := describeError(status.Error)
	logConditionalStatus[T](status, isError, value)
}

func describeError(err error) (hasError bool, description string) {
	if err != nil {
		return !errors.Is(err, context.DeadlineExceeded), err.Error()
	}
	return false, ""
}

func (c *Consumer[T]) wait() {
	c.Log.Debugf(fmt.Sprintf("wait for %d seconds...", c.Configuration.RequestDelay))
	time.Sleep(time.Duration(c.Configuration.RequestDelay) * time.Second)
}

func logConditionalStatus[T any](processStatus processor.ProcessStatus[T], isError bool, info string) {
	var (
		logConditional *zerolog.Event
		fieldName      = "info"
	)
	if !isError {
		logConditional = log.Info()
	} else {
		logConditional = log.Error()
		fieldName = "error"
	}
	if isMessageWithValues[T](processStatus.Message) {
		var message string
		if !isError {
			message = truncateText(processStatus.Message.FlatMessage)
		} else {
			message = processStatus.Message.FlatMessage
		}
		logConditional.
			Str("status", processStatus.Status).
			Str("topic", processStatus.Message.Topic).
			Str("key", processStatus.Message.Key).
			Str("message", message).
			Str(fieldName, info).
			Msg("")
	} else {
		logConditional.
			Str("status", processStatus.Status).
			Str(fieldName, info).
			Msg("")
	}
}

func isMessageWithValues[T any](message kafka.MessageForRead[T]) bool {
	return len(message.FlatMessage) > 0
}

func truncateText(text string) string {
	if len(text) > maxValueLength {
		return text[:maxValueLength] + "..."
	}
	return text
}
