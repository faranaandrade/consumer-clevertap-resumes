package updater

import "github.com/occmundial/consumer-clevertap-resumes/internal/models"

type ClevertapGetter interface {
	APICheck() error
	SendRequest(message *models.ClevetapBody) error
}
