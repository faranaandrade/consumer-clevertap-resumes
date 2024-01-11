package updater

import (
	"context"
	"errors"

	"github.com/occmundial/go-common/logger"
	"gorm.io/gorm"
)

type DataUserGetter interface {
	GetDBInfo(ctx context.Context, userID string) (string, error)
}

type DataUserRepositoryDatabase struct {
	gorm *gorm.DB
	log  logger.Logger
}

func NewDataUserRepositoryDatabase(g *gorm.DB, log *logger.Log) *DataUserRepositoryDatabase {
	return &DataUserRepositoryDatabase{
		gorm: g,
		log:  log,
	}
}

func (repository *DataUserRepositoryDatabase) GetDBInfo(ctx context.Context, userID string) (string, error) {
	var logins string
	err := repository.gorm.WithContext(ctx).
		Table("Logins").
		Select([]string{"EmailAddress"}).
		Where("LoginID = ?", userID).
		Find(&logins).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		repository.log.Error("updater", "GetEmail", err)
		return "", err
	}
	return logins, nil
}
