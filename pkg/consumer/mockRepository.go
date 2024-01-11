package consumer

import (
	"github.com/occmundial/consumer-clevertap-applies/pkg/processor"
	"github.com/stretchr/testify/mock"
)

type mockProcessor struct {
	mock.Mock
}

func (testMock *mockProcessor) CheckTargetHealth() error {
	args := testMock.Called()
	return args.Error(0)
}
func (testMock *mockProcessor) ProcessMessage() processor.ProcessStatus[int] {
	args := testMock.Called()
	return args.Get(0).(processor.ProcessStatus[int])
}
