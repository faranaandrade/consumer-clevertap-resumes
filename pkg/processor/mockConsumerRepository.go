package processor

import (
	"context"

	"github.com/occmundial/consumer-clevertap-resumes/pkg/kafka"
	"github.com/stretchr/testify/mock"
)

type mockProcessorRepository struct {
	mock.Mock
}

// CheckTargetHealth :
func (testMock *mockProcessorRepository) CheckTargetHealth() error {
	args := testMock.Called()
	return args.Error(0)
}

// GetMessage :
func (testMock *mockProcessorRepository) GetMessage() (kafka.MessageForRead[MockType], error) {
	args := testMock.Called()
	return args.Get(0).(kafka.MessageForRead[MockType]), args.Error(1)
}

func (testMock *mockProcessorRepository) CreateEvent(message any) error {
	args := testMock.Called(message)
	return args.Error(0)
}

func (testMock *mockProcessorRepository) Process(ctx context.Context, message *MockType) error {
	args := testMock.Called(ctx, message)
	return args.Error(0)
}
func (testMock *mockProcessorRepository) IsValid(message *MockType) bool {
	args := testMock.Called(message)
	return args.Bool(0)
}
func (testMock *mockProcessorRepository) GetRetryNumber(message *MockType) int {
	return message.Retries
}
func (testMock *mockProcessorRepository) SetRetryNumber(message *MockType, value int) {
	message.Retries = value
}

type MockType struct {
	Retries int
}
