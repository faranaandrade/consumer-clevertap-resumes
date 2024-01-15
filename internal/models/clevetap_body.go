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
	ResumeID int    `json:"resumeID,omitempty"`
	JobID    int    `json:"jobID,omitempty"`
	Email    string `json:"email,omitempty"`
}
