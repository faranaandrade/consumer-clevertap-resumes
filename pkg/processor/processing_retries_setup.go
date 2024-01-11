package processor

type ProcessingRetriesSetup struct {
	TopicRetry   string // Tópico para reintentos de procesamiento (opcional)
	WaitForRetry int    // Cantidad de tiempo a esperar para procesar (en segundos)
	MaxRetries   int    // Cantidad máxima de reprocesamientos permitidos en el tópico principal
	TopicMain    string
}
