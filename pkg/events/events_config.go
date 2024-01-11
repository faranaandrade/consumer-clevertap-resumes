package events

import "github.com/occmundial/consumer-clevertap-applies/config"

type Setup struct {
	APIEvents  string
	APITimeout int
	HTTPClient config.HTTPClient
}
