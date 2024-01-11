package processor

import (
	"errors"

	"github.com/occmundial/go-common/logger"

	"github.com/stretchr/testify/mock"

	"testing"
	"time"

	"github.com/occmundial/consumer-clevertap-resumes/pkg/kafka"
	l "github.com/occmundial/consumer-clevertap-resumes/pkg/location"

	"github.com/stretchr/testify/assert"
)

const MessageToSendTypeString = "MessageToSend[github.com/occmundial/consumer-clevertap-resumes/pkg/processor.MockType]"

var (
	location2        = l.GetLocation(logger.GetLogger()).Location
	topicMain        = "topic-main"
	topicRetry       = "topic-retry"
	messageToProcess = MockType{Retries: 0}
	messageForRead   = kafka.MessageForRead[MockType]{
		Key:     "c4bf18e8-10f6-48ca-aa6f-61c7e8c25578",
		Topic:   "mys-apply-successful",
		Time:    time.Date(2022, 07, 11, 10, 0, 0, 0, location2),
		Message: messageToProcess,
	}
	getLogger         = logger.GetLogger()
	MessageToSendType = mock.AnythingOfType(MessageToSendTypeString)
)

func newMockProcessor(mockRepository *mockProcessorRepository, topicRetry string, maxRetries, waitForRetry int) Processor[MockType] {
	return Processor[MockType]{
		QueueGetter:     mockRepository,
		Retryer:         mockRepository,
		MessageProcesor: mockRepository,
		Location:        l.GetLocation(getLogger),
		Configuration: &ProcessingRetriesSetup{
			TopicRetry:   topicRetry,
			MaxRetries:   maxRetries,
			WaitForRetry: waitForRetry,
			TopicMain:    topicMain},
		Log: getLogger,
	}
}

func TestConsumerService_When_RetryNumberIsLowerThanMaxRetries_Expect_RetryInMainTopicAndIncrementsRetryNumber(t *testing.T) {
	mockRepository := new(mockProcessorRepository)
	messageToProcess := MockType{Retries: 2}
	processor := newMockProcessor(mockRepository, topicRetry, 3, 30)
	nextTopic := processor.getRetryTopic(&messageToProcess)
	assert.Equal(t, topicMain, nextTopic)
	assert.Equal(t, 3, messageToProcess.Retries)
}

func TestConsumerService_When_RetryNumberIsMaxRetries_Expect_RetryInRetryTopicAndInitializeRetryNumber(t *testing.T) {
	mockRepository := new(mockProcessorRepository)
	messageToProcess := MockType{Retries: 3}
	processor := newMockProcessor(mockRepository, topicRetry, 3, 30)
	nextTopic := processor.getRetryTopic(&messageToProcess)
	assert.Equal(t, topicRetry, nextTopic)
	assert.Equal(t, 0, messageToProcess.Retries)
}

func TestConsumerService_When_MaxRetriesIsZeroAndNotTopicRetry_Expect_NotRetryTopic(t *testing.T) {
	mockRepository := new(mockProcessorRepository)
	messageToProcess := MockType{Retries: 0}
	processor := newMockProcessor(mockRepository, "", 0, 30)
	nextTopic := processor.getRetryTopic(&messageToProcess)
	assert.Equal(t, "", nextTopic)
	assert.Equal(t, 0, messageToProcess.Retries)
}

func TestConsumerService_When_WaitForRetryGreaterZeroAndConsumeCompleteTime_Expect_NoWaitsTime(t *testing.T) {
	testWaitTimes(t, 20, 50, 0)
}

func TestConsumerService_When_WaitForRetryGreaterZeroAndNotConsumeTime_Expect_WaitsCompleteTime(t *testing.T) {
	testWaitTimes(t, 20, 20, 30)
}

func TestConsumerService_When_WaitForRetryGreaterZeroAndConsumeTime_Expect_WaitsPartiallyTime(t *testing.T) {
	testWaitTimes(t, 20, 30, 20)
}

func testWaitTimes(t *testing.T, firstSec, secondSec int, expectedWaitingTime float64) {
	mockRepository := new(mockProcessorRepository)
	mockRepository.On("CheckTargetHealth").Return(nil)
	mockRepository.On("GetMessage").Return(messageForRead, nil)
	consumerService := newMockProcessor(mockRepository, "", 0, 30)
	messageTime := time.Date(2022, 7, 11, 3, 40, firstSec, 0, location2)
	currentTime := time.Date(2022, 7, 11, 3, 40, secondSec, 0, location2)
	actualWaitingTime := consumerService.calculateWaitingTime(messageTime, currentTime).Seconds()
	assert.Equal(t, expectedWaitingTime, actualWaitingTime)
}

// todo perfecto => happy path
func TestConsumerService_When_MessageIsOK_Expect_StatusProcessOK(t *testing.T) {
	mockRepository := new(mockProcessorRepository)
	mockRepository.On("CheckTargetHealth").Return(nil)
	mockRepository.On("GetMessage").Return(messageForRead, nil)
	mockRepository.On("IsValid", &messageForRead.Message).Return(true)
	mockRepository.On("Process", mock.Anything, &messageForRead.Message).Return(nil)
	expectedProcessStatus := ProcessStatus[MockType]{
		Status:  StatusFullProcessOK,
		Message: messageForRead,
	}
	processor := newMockProcessor(mockRepository, "", 0, 30)
	actualProcessStatus := processor.ProcessMessage()
	assert.Equal(t, expectedProcessStatus, actualProcessStatus)
}

// Error en el procesamiento del mensaje
func TestConsumerService_When_MessageProcessedError_Expect_StatusProcessError(t *testing.T) {
	mockRepository := new(mockProcessorRepository)

	err := errors.New("error en procesamiento del mensaje")
	mockRepository.On("CheckTargetHealth").Return(nil)
	mockRepository.On("GetMessage").Return(messageForRead, nil)
	mockRepository.On("IsValid", &messageForRead.Message).Return(true)
	mockRepository.On("Process", mock.Anything, &messageForRead.Message).Return(err)

	expectedProcessStatus := ProcessStatus[MockType]{
		Status:  StatusProcessError,
		Message: messageForRead,
		Error:   err,
	}

	processor := newMockProcessor(mockRepository, "", 0, 0)
	actualProcessStatus := processor.ProcessMessage()

	assert.Equal(t, expectedProcessStatus, actualProcessStatus)
}

// Error durante la lectura del mensaje a procesar
func TestConsumerService_When_GetMessageError_Expect_StatusReadMessageError(t *testing.T) {
	mockRepository := new(mockProcessorRepository)

	err := errors.New("error durante la lectura del mensaje a procesar")
	mockRepository.On("CheckTargetHealth").Return(nil)
	mockRepository.On("GetMessage").Return(messageForRead, err)

	expectedProcessStatus := ProcessStatus[MockType]{
		Status:  StatusReadMessageError,
		Message: messageForRead,
		Error:   err,
	}

	processor := newMockProcessor(mockRepository, "", 0, 0)
	actualProcessStatus := processor.ProcessMessage()

	assert.Equal(t, expectedProcessStatus, actualProcessStatus)
}

// Error el mensaje a procesar no es válido
func TestConsumerService_When_InvalidMessage_Expect_StatusInvalidMessage(t *testing.T) {
	mockRepository := new(mockProcessorRepository)

	mockRepository.On("CheckTargetHealth").Return(nil)
	msgForRead := kafka.MessageForRead[MockType]{}
	mockRepository.On("GetMessage").Return(msgForRead, nil)
	mockRepository.On("IsValid", &messageForRead.Message).Return(false)

	expectedProcessStatus := ProcessStatus[MockType]{
		Status:  StatusInvalidMessage,
		Message: msgForRead,
		Error:   nil,
	}

	processor := newMockProcessor(mockRepository, "", 0, 0)
	actualProcessStatus := processor.ProcessMessage()

	assert.Equal(t, expectedProcessStatus, actualProcessStatus)
}

// Error el mensaje a procesar no es válido
func TestConsumerService_When_WaitingTime(t *testing.T) {
	mockRepository := new(mockProcessorRepository)

	mockRepository.On("CheckTargetHealth").Return(nil)
	msgForRead := kafka.MessageForRead[MockType]{Time: time.Now().Add(time.Second * 1)}
	mockRepository.On("GetMessage").Return(msgForRead, nil)
	mockRepository.On("IsValid", &messageForRead.Message).Return(true)
	mockRepository.On("Process", mock.Anything, &messageForRead.Message).Return(nil)

	expectedProcessStatus := ProcessStatus[MockType]{
		Status:  StatusFullProcessOK,
		Message: msgForRead,
		Error:   nil,
	}

	processor := newMockProcessor(mockRepository, "", 0, 1)
	actualProcessStatus := processor.ProcessMessage()

	assert.Equal(t, expectedProcessStatus, actualProcessStatus)
}

func TestConsumerService_successful_sendRetryMessage(t *testing.T) {
	mockRepository := new(mockProcessorRepository)
	err := errors.New("error en procesamiento del mensaje")
	mockRepository.On("GetMessage").Return(messageForRead, nil)
	mockRepository.On("IsValid", &messageForRead.Message).Return(true)
	mockRepository.On("Process", mock.Anything, &messageForRead.Message).Return(err)
	mockRepository.On("CheckTargetHealth").Return(nil)
	mockRepository.On("CreateEvent", MessageToSendType).Return(nil)

	expectedProcessStatus := ProcessStatus[MockType]{
		Status:  StatusProcessError,
		Message: messageForRead,
		Error:   err,
	}

	processor := newMockProcessor(mockRepository, "topicRetry", 0, 0)
	actualProcessStatus := processor.ProcessMessage()

	assert.Equal(t, expectedProcessStatus, actualProcessStatus)
}

func TestConsumerService_error_sendRetryMessage(t *testing.T) {
	mockRepository := new(mockProcessorRepository)
	err := errors.New("error en procesamiento del mensaje")
	mockRepository.On("CheckTargetHealth").Return(nil)
	mockRepository.On("GetMessage").Return(messageForRead, nil)
	mockRepository.On("Process", mock.Anything, &messageForRead.Message).Return(err)
	mockRepository.On("IsValid", &messageForRead.Message).Return(true)
	mockRepository.On("CreateEvent", MessageToSendType).Return(err)

	expectedProcessStatus := ProcessStatus[MockType]{
		Status:  StatusProcessError,
		Message: messageForRead,
		Error:   err,
	}

	processor := newMockProcessor(mockRepository, "topicRetry", 0, 0)
	actualProcessStatus := processor.ProcessMessage()

	assert.Equal(t, expectedProcessStatus, actualProcessStatus)
}

func TestConsumerService_CheckTargetHealth(t *testing.T) {
	mockRepository := new(mockProcessorRepository)
	mockRepository.On("CheckTargetHealth").Return(nil)
	processor := newMockProcessor(mockRepository, "topicRetry", 0, 0)
	assert.NoError(t, processor.CheckTargetHealth())
}
