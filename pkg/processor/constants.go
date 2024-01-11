package processor

const (
	// StatusReadMessageError error durante la lectura del mensaje
	StatusReadMessageError = "read-message-error"
	// StatusInvalidMessage el mensaje no tiene el mínimo de valores para considerarlo como válido
	StatusInvalidMessage = "invalid-message-error"
	// StatusFullProcessOK Procesamiento completo del mensaje (mensaje procesado y borrado de la cola)
	StatusFullProcessOK = "full-processed-message"
	// StatusProcessError Error en procesamiento del mensaje: falló el procesamiento del mensaje
	StatusProcessError = "processing-error"
)
