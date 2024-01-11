package test

import (
	"context"

	"github.com/occmundial/consumer-clevertap-applies/internal/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (a *RepositoryMock) APICheck() error {
	args := a.Mock.Called()
	return args.Error(0)
}

func (a *RepositoryMock) SendRequest(message *models.ClevetapBody) error {
	args := a.Mock.Called(message)
	return args.Error(0)
}

func (a *RepositoryMock) GetDBInfo(ctx context.Context, userID string) (string, error) {
	args := a.Mock.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
