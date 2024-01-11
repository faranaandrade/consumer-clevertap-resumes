package updater

import (
	"context"
	"errors"
	"testing"

	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
	"github.com/occmundial/consumer-clevertap-resumes/test"
	"github.com/occmundial/go-common/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const address = "test@email.com"

func TestIsValid(t *testing.T) {
	resolver := NewMessageResolverUpdater(&APIClevetapRepository{}, &DataUserRepositoryDatabase{}, logger.GetLogger())
	validMessage := &models.MessageToProcess{
		ResumeID: 1234,
		JobID:    123456,
		UserID:   "123456",
	}
	invalidMessage := &models.MessageToProcess{
		ResumeID: 0,
		JobID:    0,
		UserID:   "",
	}
	if !resolver.IsValid(validMessage) {
		t.Error("Se esperaba que el mensaje válido sea verdadero")
	}
	if resolver.IsValid(invalidMessage) {
		t.Error("Se esperaba que el mensaje inválido sea falso")
	}
}

func TestGetAndSetRetryNumber(t *testing.T) {
	resolver := NewMessageResolverUpdater(&APIClevetapRepository{}, &DataUserRepositoryDatabase{}, logger.GetLogger())
	message := &models.MessageToProcess{RetryNumber: 0}
	value := 33
	resolver.SetRetryNumber(message, value)
	if resolver.GetRetryNumber(message) != value {
		t.Error("Se esperaba que el valor sea modificado")
	}
}

func TestProcessSuccessful(t *testing.T) {
	cxt := context.Background()
	userID := "User"
	messageToProcess := models.MessageToProcess{RetryNumber: 0, CreationDate: "2006-01-02T15:04:05Z07:00", UserID: userID}
	mockRepository := new(test.RepositoryMock)
	mockRepository.On("APICheck").Return(nil)
	mockRepository.On("SendRequest", mock.Anything).Return(nil)
	mockRepository.On("GetDBInfo", cxt, userID).Return(address, nil)
	resolver := MessageResolverUpdater{PageSize: 1, ClevetapRepository: mockRepository,
		DataUserRepository: mockRepository, log: logger.GetLogger()}
	message := &messageToProcess
	assert.NoError(t, resolver.Process(cxt, message))
	mockRepository.AssertExpectations(t)
}

func TestProcessFailSendRequest(t *testing.T) {
	cxt := context.Background()
	userID := "userTest"
	process := models.MessageToProcess{RetryNumber: 0, CreationDate: "2006-01-02T15:04:05Z07:00", UserID: userID}
	mockRepository := new(test.RepositoryMock)
	mockRepository.On("APICheck").Return(nil)
	mockRepository.On("SendRequest", mock.Anything).Return(errors.New("fakeError"))
	mockRepository.On("GetDBInfo", cxt, userID).Return(address, nil)
	resolver := MessageResolverUpdater{PageSize: 1, ClevetapRepository: mockRepository,
		DataUserRepository: mockRepository, log: logger.GetLogger()}
	assert.Error(t, resolver.Process(cxt, &process))
	mockRepository.AssertExpectations(t)
}

func Test_GetEmail_fail_then_return_error(t *testing.T) {
	cxt := context.Background()
	userID := "userId"
	process := models.MessageToProcess{RetryNumber: 0, CreationDate: "2006-01-02T15:04:05Z07:00", UserID: userID}
	mockRepository := new(test.RepositoryMock)
	mockRepository.On("APICheck").Return(nil)
	mockRepository.On("GetDBInfo", cxt, userID).Return("", errors.New("fakeError"))
	resolver := MessageResolverUpdater{PageSize: 1, ClevetapRepository: mockRepository,
		DataUserRepository: mockRepository, log: logger.GetLogger()}
	assert.Error(t, resolver.Process(cxt, &process))
	mockRepository.AssertExpectations(t)
}

func TestProcessFailCheckHealth(t *testing.T) {
	mockRepository := new(test.RepositoryMock)
	mockRepository.On("APICheck").Return(errors.New("fakeError"))
	resolver := MessageResolverUpdater{PageSize: 1, ClevetapRepository: mockRepository}
	assert.Error(t, resolver.Process(context.Background(), &models.MessageToProcess{RetryNumber: 0}))
	mockRepository.AssertExpectations(t)
}

func TestMapEmailStatisticForCreation(t *testing.T) {
	resolver := NewMessageResolverUpdater(&APIClevetapRepository{}, &DataUserRepositoryDatabase{}, logger.GetLogger())
	result := resolver.mapEmailStatisticForCreation(12, 12, "test@gmail.com")
	assert.True(t, len(result.Body) > 0)
	assert.Equal(t, "test@gmail.com", result.Body[0].EvtData.Email)
	assert.Equal(t, 12, result.Body[0].EvtData.JobID)
}
