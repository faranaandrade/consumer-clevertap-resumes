package events

import (
	"context"
	"errors"
	"fmt"

	"net/http"
	"strings"
	"time"

	httpclient "github.com/occmundial/go-common/http/client"
)

// ProcessAPIHealth : regresa el statusCode del endpoint de salud
func ProcessAPIHealth(httpClient *httpclient.StandardClient, timeout int, urlAPI string, chanErrorMsg chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	requestOptions := httpclient.NewGetRequestOptions(urlAPI)
	response, err := httpClient.Request(ctx, requestOptions)
	if err == nil {
		defer response.Body.Close()
		if response.StatusCode == http.StatusOK {
			chanErrorMsg <- ""
		} else {
			chanErrorMsg <- fmt.Sprintf("Processor Unavailable -> %s", urlAPI)
		}
	} else {
		chanErrorMsg <- fmt.Sprintf("Processor Unavailable -> %s", err.Error())
	}
}

func ConcatErrors(apiErrors ...string) error {
	errorMessages := ""
	for _, msg := range apiErrors {
		if len(msg) > 0 {
			errorMessages += msg + "\n"
		}
	}
	if len(errorMessages) > 0 {
		return errors.New(strings.TrimSuffix(errorMessages, "\n"))
	}
	return nil
}

func CloseChannels(channels ...chan string) {
	for _, item := range channels {
		close(item)
	}
}
