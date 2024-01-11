package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/occmundial/go-common/http/body"

	"github.com/occmundial/consumer-clevertap-applies/pkg/location"
	httpclient "github.com/occmundial/go-common/http/client"

	"net/http"
	"strings"
	"time"
)

var (
	urlMessages     string
	urlEventsHealth string
)

// NewEventsRepository : Factory que crea un "EventsRepository"
func NewEventsRepository[T any](configuration *Setup, locater *location.Location) *Repository[T] {
	standardClientSetup := httpclient.StandardClientSetup{APITimeout: configuration.HTTPClient.APITimeout}
	retryableclientSetup := httpclient.RetryableClientSetup(configuration.HTTPClient)
	cr := Repository[T]{
		EventsSetup:         configuration,
		Locater:             locater,
		HTTPStandardClient:  httpclient.NewStandardClient(&standardClientSetup),
		HTTPRetryableClient: httpclient.NewRetryableClient(&retryableclientSetup),
	}
	cr.init()
	return &cr
}

// EventsRepository :
type Repository[T any] struct {
	EventsSetup         *Setup
	Locater             location.Locater
	HTTPStandardClient  *httpclient.StandardClient
	HTTPRetryableClient *httpclient.RetryableClient
}

func (r Repository[T]) init() {
	urlMessages = fmt.Sprintf("%s/events/v2/messages", strings.TrimSuffix(r.EventsSetup.APIEvents, "/"))
	urlEventsHealth = fmt.Sprintf("%s/events/health", strings.TrimSuffix(r.EventsSetup.APIEvents, "/"))
}

// CheckTargetHealth :
func (r Repository[T]) CheckTargetHealth() error {
	// declaración de canales para el chequeo de salud de dependencies
	chanEventsHealth := make(chan string)
	// cerramos los canales de manera diferida
	defer CloseChannels(chanEventsHealth)
	// llamadas asíncronas a las funciones de chequeo
	go ProcessAPIHealth(r.HTTPStandardClient, r.EventsSetup.APITimeout, urlEventsHealth, chanEventsHealth)
	eventsHealth := <-chanEventsHealth
	return ConcatErrors(eventsHealth)
}

func (r *Repository[T]) CreateEvent(message any) error {
	jsonValue, err := json.Marshal(message)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.EventsSetup.APITimeout)*time.Second)
	defer cancel()
	requestOptions := httpclient.NewPostRequestOptions(urlMessages, "application/json", bytes.NewBuffer(jsonValue))
	response, err := r.HTTPRetryableClient.Request(ctx, requestOptions)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err == nil && response.StatusCode != http.StatusCreated {
		return body.ResponseToError(response)
	}
	return err
}
