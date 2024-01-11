package models

// MessageToProcess : mensaje a procesar
type MessageToProcess struct {
	RetryNumber  int    `json:"retryNumber"` // Iteración actual de reprocesamiento en el tópico actual (main topic)
	UserID       string `json:"userid"`
	JobID        int    `json:"jobId"`
	CreationDate string `json:"creationDate"`
	ResumeID     int    `json:"resumeId"`
}
