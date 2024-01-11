package models

type ResponseClevertap struct {
	Status      string `json:"status"`
	Processed   int    `json:"processed"`
	Unprocessed []any  `json:"unprocessed"`
}
