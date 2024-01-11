package setup

import (
	"github.com/occmundial/consumer-clevertap-resumes/config"
	"github.com/occmundial/consumer-clevertap-resumes/pkg/gorm"
)

func NewGormSetup(configuration *config.Configuration) *gorm.Setup {
	return &gorm.Setup{
		User:               configuration.GormConnection.User,
		Password:           configuration.GormConnection.Password,
		Server:             configuration.GormConnection.Server,
		Database:           configuration.GormConnection.Database,
		ConnectTimeOut:     configuration.GormConnection.ConnectTimeOut,
		Environment:        configuration.Environment,
		MaxIdleConnections: configuration.GormConnection.MaxIdleConnections,
		MaxOpenConnections: configuration.GormConnection.MaxOpenConnections,
		ConnMaxLifetime:    configuration.GormConnection.ConnMaxLifetime,
		Port:               configuration.GormConnection.Port,
	}
}
