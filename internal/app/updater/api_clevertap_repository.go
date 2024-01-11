package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/occmundial/consumer-clevertap-applies/config"
	"github.com/occmundial/consumer-clevertap-applies/internal/models"
	httpbody "github.com/occmundial/go-common/http/body"
	httpclient "github.com/occmundial/go-common/http/client"
)

const (
	success            = "success"
	accountIDHeaderKey = "X-CleverTap-Account-Id"
	ClevertapHeaderkey = "X-CleverTap-Passcode"
	contentType        = "application/json"
	limit              = 1048576
)

type APIClevetapRepository struct {
	ClevertapAPISetup   *config.ClevertapAPISetup
	HTTPStandardClient  *httpclient.StandardClient
	HTTPRetryableClient *httpclient.RetryableClient
}

func NewAPIClevetapRepository(configuration *config.Configuration) *APIClevetapRepository {
	standardClientSetup := httpclient.StandardClientSetup{APITimeout: configuration.ClevertapConfig.HTTPClientSetup.APITimeout}
	retryableclientSetup := httpclient.RetryableClientSetup(configuration.ClevertapConfig.HTTPClientSetup)
	return &APIClevetapRepository{
		ClevertapAPISetup:   &configuration.ClevertapConfig.APISetup,
		HTTPStandardClient:  httpclient.NewStandardClient(&standardClientSetup),
		HTTPRetryableClient: httpclient.NewRetryableClient(&retryableclientSetup),
	}
}

// APICheck :
func (r *APIClevetapRepository) APICheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.ClevertapAPISetup.Timeout)*time.Second)
	defer cancel()
	url := fmt.Sprintf("%s/health", strings.TrimSuffix(r.ClevertapAPISetup.Host, "/"))
	response, err := r.HTTPStandardClient.Request(ctx, httpclient.NewGetRequestOptions(url))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return httpbody.ResponseToError(response)
	}
	defer response.Body.Close()
	return nil
}

func (r *APIClevetapRepository) SendRequest(message *models.ClevetapBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.ClevertapAPISetup.Timeout)*time.Second)
	defer cancel()
	msgBody, err := httpbody.ModelToBody(message)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/1/upload", strings.TrimSuffix(r.ClevertapAPISetup.Host, "/"))
	requestOptions := httpclient.NewPostRequestOptions(url, contentType, msgBody).
		WithHeader(accountIDHeaderKey, r.ClevertapAPISetup.AccountID).
		WithHeader(ClevertapHeaderkey, r.ClevertapAPISetup.ClevertapPasscode)
	response, err := r.HTTPRetryableClient.Request(ctx, requestOptions)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return httpbody.ResponseToError(response)
	}
	defer response.Body.Close()
	body, err := r.getBody(response)
	if err != nil {
		return err
	}
	responseClevertap := models.ResponseClevertap{}
	err = json.Unmarshal(body, &responseClevertap)
	if err != nil {
		return err
	}
	if responseClevertap.Status != success || len(responseClevertap.Unprocessed) > 0 {
		errMessage := fmt.Sprintf("api.response.StatusCode = %s", http.StatusText(response.StatusCode))
		return fmt.Errorf("%s | api.response.Body = %s", errMessage, string(body))
	}
	return nil
}

// getBody
func (r *APIClevetapRepository) getBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	// LimitReader: protect against malicious attacks on your server
	return io.ReadAll(io.LimitReader(response.Body, limit))
}
