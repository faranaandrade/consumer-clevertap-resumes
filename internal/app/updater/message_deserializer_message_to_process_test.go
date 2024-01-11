package updater

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
)

func TestGetMessageFromBytes(t *testing.T) {
	deserializer := NewMessageDeserializer()
	mockMessage := models.MessageToProcess{
		RetryNumber: 1,
		UserID:      "TestUser",
	}
	mockBytes, _ := json.Marshal(mockMessage)
	result, err := deserializer.GetMessageFromBytes(mockBytes)
	if err != nil {
		t.Errorf("Se produjo un error inesperado: %v", err)
	}
	if !reflect.DeepEqual(result, mockMessage) {
		t.Errorf("Se esperaba %v pero se obtuvo %v", mockMessage, result)
	}
}
