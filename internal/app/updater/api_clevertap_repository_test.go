package updater

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/occmundial/consumer-clevertap-applies/config"
	"github.com/occmundial/consumer-clevertap-applies/internal/models"
)

const (
	URL      = "https://example.com"
	HeathURL = URL + "/health"
)

func TestAPICheck_with_OKStatus(t *testing.T) {
	setup := config.ClevertapAPISetup{
		Host:    URL,
		Timeout: 5,
	}
	configuration := config.Configuration{ClevertapConfig: config.ClevertapConfig{APISetup: setup}}
	repository := NewAPIClevetapRepository(&configuration)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, HeathURL,
		httpmock.NewStringResponder(http.StatusOK, ``))

	err := repository.APICheck()
	if err != nil {
		t.Errorf("Se produjo un error inesperado en APICheck: %v", err)
	}
}

func TestAPICheck_with_FailStatus(t *testing.T) {
	setup := config.ClevertapAPISetup{
		Host:    URL,
		Timeout: 5,
	}
	configuration := config.Configuration{ClevertapConfig: config.ClevertapConfig{APISetup: setup}}
	repository := NewAPIClevetapRepository(&configuration)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, HeathURL,
		httpmock.NewStringResponder(http.StatusInternalServerError, ``))

	err := repository.APICheck()
	if err == nil {
		t.Errorf("Se esperaba un error en APICheck: %v", err)
	}
}

func TestAPICheck_with_error_request(t *testing.T) {
	setup := config.ClevertapAPISetup{
		Host:    URL,
		Timeout: 5,
	}
	configuration := config.Configuration{ClevertapConfig: config.ClevertapConfig{APISetup: setup}}
	repository := NewAPIClevetapRepository(&configuration)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, HeathURL,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("internal server error")
		})

	err := repository.APICheck()
	if err == nil {
		t.Errorf("Se esperaba un error en APICheck: %v", err)
	}
}

func TestSendRequest(t *testing.T) {
	err := sendMockRequest(t, http.StatusOK, models.ResponseClevertap{Status: success})
	if err != nil {
		t.Errorf("Se produjo un error inesperado en SendRequest: %v", err)
	}
}

func TestSendRequest_with_error_code_status(t *testing.T) {
	err := sendMockRequest(t, http.StatusNotFound, models.ResponseClevertap{Status: success})
	if err == nil {
		t.Errorf("Se esperaba un error en SendRequest: %v", err)
	}
}

func TestSendRequest_with_fail_status(t *testing.T) {
	err := sendMockRequest(t, http.StatusOK, models.ResponseClevertap{Status: "fail"})
	if err == nil {
		t.Errorf("Se esperaba un error en SendRequest: %v", err)
	}
}

func sendMockRequest(t *testing.T, status int, rc models.ResponseClevertap) error {
	configuration := config.Configuration{ClevertapConfig: config.ClevertapConfig{APISetup: config.ClevertapAPISetup{
		Timeout: 5,
	}}}
	repository := NewAPIClevetapRepository(&configuration)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("bad method: %s", r.Method)
		}
		w.WriteHeader(status)
		body, err := json.Marshal(rc)
		if err != nil {
			t.Fatalf("error writing")
		}
		_, err = w.Write(body)
		if err != nil {
			t.Fatalf("error writing")
		}
	}))
	defer ts.Close()
	configuration.ClevertapConfig.APISetup.Host = ts.URL
	err := repository.SendRequest(&models.ClevetapBody{})
	return err
}

func TestSendRequest_with_error_request(t *testing.T) {
	configuration := config.Configuration{ClevertapConfig: config.ClevertapConfig{APISetup: config.ClevertapAPISetup{
		Timeout: 5,
	}}}
	repository := NewAPIClevetapRepository(&configuration)
	err := repository.SendRequest(nil)
	if err == nil {
		t.Errorf("Se esperaba un error en SendRequest: %v", err)
	}
}
