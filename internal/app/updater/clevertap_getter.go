package updater

import "github.com/occmundial/consumer-clevertap-applies/internal/models"

type ClevertapGetter interface {
	APICheck() error
	SendRequest(message *models.ClevetapBody) error
}
