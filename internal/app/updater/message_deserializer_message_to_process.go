package updater

import (
	"encoding/json"

	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
)

func NewMessageDeserializer() *MessageDeserializerMessageToProcess {
	return &MessageDeserializerMessageToProcess{}
}

type MessageDeserializerMessageToProcess struct {
}

func (h *MessageDeserializerMessageToProcess) GetMessageFromBytes(flatMessage []byte) (models.MessageToProcess, error) {
	msgToProcess := models.MessageToProcess{}
	err := json.Unmarshal(flatMessage, &msgToProcess)
	return msgToProcess, err
}
