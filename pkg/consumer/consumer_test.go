package consumer

import (
	"errors"
	"testing"

	"github.com/occmundial/go-common/logger"

	"github.com/occmundial/consumer-clevertap-applies/pkg/kafka"
	"github.com/occmundial/consumer-clevertap-applies/pkg/processor"
)

func TestIsTargetHealth_when_is_good_return_true(t *testing.T) {
	mocked := new(mockProcessor)
	mocked.On("CheckTargetHealth").Return(nil)
	consumer := Consumer[int]{
		Processor: mocked,
		Log:       logger.GetLogger(),
	}
	if !consumer.isTargetHealth() {
		t.Error("Esperaba que la salud sea verdadera")
	}
}

func TestIsTargetHealth_when_is_bad_return_false(t *testing.T) {
	mocked := new(mockProcessor)
	err := errors.New("testError")
	mocked.On("CheckTargetHealth").Return(err)
	consumer := Consumer[int]{
		Processor: mocked,
		Log:       logger.GetLogger(),
	}
	if consumer.isTargetHealth() {
		t.Error("Esperaba que la salud sea falsa")
	}
}

func TestProcessMessage_when_is_ok_return_true(t *testing.T) {
	mocked := new(mockProcessor)
	mocked.On("ProcessMessage").Return(processor.ProcessStatus[int]{})
	consumer := Consumer[int]{
		Processor: mocked,
		Log:       logger.GetLogger(),
	}
	if !consumer.processMessage() {
		t.Error("Esperaba que el procesamiento del mensaje sea verdadero")
	}
}

func TestProcessMessage_when_is_bad_return_false(t *testing.T) {
	mocked := new(mockProcessor)
	message := kafka.MessageForRead[int]{
		FlatMessage: "valor",
		Message:     3,
	}
	arguments := processor.ProcessStatus[int]{Status: StatusProcessStartError, Message: message}
	mocked.On("ProcessMessage").Return(arguments)
	consumer := Consumer[int]{
		Processor: mocked,
		Log:       logger.GetLogger(),
	}
	if consumer.processMessage() {
		t.Error("Esperaba que el procesamiento del mensaje sea false")
	}
}

func TestProcessMessage_when_is_error_process_return_false(t *testing.T) {
	mocked := new(mockProcessor)
	message := kafka.MessageForRead[int]{
		FlatMessage: "valor muy muy muy muy muy muy muy muy muy muy muy muy muy muy muy muy largo para ser leído en los logs.",
		Message:     3,
	}
	arguments := processor.ProcessStatus[int]{Status: StatusProcessStartError, Message: message}
	mocked.On("ProcessMessage").Return(arguments)
	consumer := Consumer[int]{
		Processor: mocked,
		Log:       logger.GetLogger(),
	}
	if consumer.processMessage() {
		t.Error("Esperaba que el procesamiento del mensaje sea false")
	}
}

func TestProcessMessage_empty_message_when_is_bad_return_false(t *testing.T) {
	mocked := new(mockProcessor)
	message := kafka.MessageForRead[int]{
		FlatMessage: "mockValor",
		Message:     3,
	}
	arguments := processor.ProcessStatus[int]{Error: errors.New("fakeError"), Message: message}
	mocked.On("ProcessMessage").Return(arguments)
	consumer := Consumer[int]{
		Processor: mocked,
		Log:       logger.GetLogger(),
	}
	if consumer.processMessage() {
		t.Error("Esperaba que el procesamiento del mensaje sea false")
	}
}

func TestConsumeWithErrors(t *testing.T) {
	mocked := new(mockProcessor)
	mocked.On("CheckTargetHealth").Return(nil)
	mocked.On("ProcessMessage").Return(processor.ProcessStatus[int]{Error: errors.New("fakeError")})
	consumer := Consumer[int]{
		Processor:     mocked,
		Log:           logger.GetLogger(),
		Configuration: &Setup{ArtifactVersion: "1"},
	}
	consumer.consume()
	mocked.AssertExpectations(t)
}

func TestConsumeWithBadHealth(t *testing.T) {
	mocked := new(mockProcessor)
	mocked.On("CheckTargetHealth").Return(errors.New("fakeError"))
	consumer := Consumer[int]{
		Processor:     mocked,
		Log:           logger.GetLogger(),
		Configuration: &Setup{ArtifactVersion: "1"},
	}
	consumer.consume()
	mocked.AssertExpectations(t)
}
