package models

// MessageToProcess : mensaje a procesar
type MessageToProcess struct {
	RetryNumber       int    `json:"retryNumber"` // Iteración actual de reprocesamiento en el tópico actual (main topic)
	UserID            string `json:"userid"`
	JobID             int    `json:"jobId"`
	CreationDate      string `json:"creationDate"`
	ResumeID          int    `json:"resumeId"`
	CvReady           bool   `json:"cv_ready"`
	CvAttached        bool   `json:"cv_attached"`
	EducationLevel    string `json:"education_level"`
	YearsOfExperience int    `json:"years_of_experience"`
}
