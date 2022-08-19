package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"github.com/fingerprintjs/fingerprint-pro-server-api-go-sdk/sdk"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func getMockResponse(path string) sdk.Response {
	var mockResponse sdk.Response

	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &mockResponse)

	if err != nil {
		log.Fatal(err)
	}

	return mockResponse
}

func readConfig() map[string]interface{} {
	fileName := "../config.json"
	configContents, err := os.ReadFile(fileName)

	var config map[string]interface{}

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(configContents, &config); err != nil {
		log.Fatal(err)
	}

	configContents, err = json.MarshalIndent(config, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func TestReturnsVisits(t *testing.T) {
	mockResponse := getMockResponse("../test/mocks/visits_limit_1.json")

	ts := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		config := readConfig()
		integrationInfo := r.URL.Query().Get("ii")
		assert.Equal(t, integrationInfo, fmt.Sprintf("fingerprint-pro-server-go-sdk/%s", config["packageVersion"]))

		apiKey := r.Header.Get("Auth-Api-Key")
		assert.Equal(t, apiKey, "api_key")

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(mockResponse)

		if err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()

	cfg := sdk.NewConfiguration()
	cfg.ChangeBasePath(ts.URL)

	client := sdk.NewAPIClient(cfg)

	ctx := context.WithValue(context.Background(), sdk.ContextAPIKey, sdk.APIKey{
		Key: "api_key",
	})

	res, _, err := client.FingerprintApi.GetVisits(ctx, "visitor_id", nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.VisitorId, mockResponse.VisitorId)
}

func TestReturnsVisitsWithPagination(t *testing.T) {
	config := sdk.FingerprintApiGetVisitsOpts{
		RequestId: optional.NewString("request_id"),
		Before:    optional.NewInt32(10),
		Limit:     optional.NewInt32(500),
		LinkedId:  optional.NewString("request_id"),
	}

	mockResponse := getMockResponse("../test/mocks/visits_limit_500.json")

	ts := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		parseErr := r.ParseForm()

		assert.NoError(t, parseErr)

		assert.Equal(t, r.Form.Get("request_id"), config.RequestId.Value())
		assert.Equal(t, r.Form.Get("before"), fmt.Sprint(config.Before.Value()))
		assert.Equal(t, r.Form.Get("limit"), fmt.Sprint(config.Limit.Value()))
		assert.Equal(t, r.Form.Get("linked_id"), config.LinkedId.Value())

		apiKey := r.Header.Get("Auth-Api-Key")
		assert.Equal(t, apiKey, "api_key")

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(mockResponse)

		if err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()

	cfg := sdk.NewConfiguration()
	cfg.ChangeBasePath(ts.URL)

	client := sdk.NewAPIClient(cfg)

	ctx := context.WithValue(context.Background(), sdk.ContextAPIKey, sdk.APIKey{
		Key: "api_key",
	})

	res, _, err := client.FingerprintApi.GetVisits(ctx, "visitor_id", &config)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.VisitorId, mockResponse.VisitorId)
	assert.Equal(t, res.Visits, mockResponse.Visits)
	assert.Equal(t, res.PaginationKey, mockResponse.PaginationKey)
	assert.Equal(t, res.LastTimestamp, mockResponse.LastTimestamp)
}
