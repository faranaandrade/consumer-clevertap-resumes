package events

import "github.com/occmundial/consumer-clevertap-resumes/config"

type Setup struct {
	APIEvents  string
	APITimeout int
	HTTPClient config.HTTPClient
}
