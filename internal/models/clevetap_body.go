package models

// ClevetapBody :
type ClevetapBody struct {
	Body []ClevertapData `json:"d"`
}

// ClevertapData :
type ClevertapData struct {
	Identity string  `json:"identity"`
	TS       int64   `json:"ts,omitempty"`
	TypeUse  string  `json:"type"`
	EvtName  string  `json:"evtName"`
	EvtData  EvtData `json:"evtData"`
}

// EvtData :
type EvtData struct {
	ResumeID          int    `json:"resumeID,omitempty"`
	Email             string `json:"email"`
	CvReady           bool   `json:"cv_ready,omitempty"`
	CvAttached        bool   `json:"cv_attached,omitempty"`
	EducationLevel    string `json:"education_level,omitempty"`
	YearsOfExperience int    `json:"years_of_experience,omitempty"`
}
