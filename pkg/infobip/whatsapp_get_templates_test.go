package infobip

import (
	"context"
	"encoding/json"
	"fmt"
	"infobip-go-client/pkg/infobip/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTemplatesValidReq(t *testing.T) {
	apiKey := "secret"
	sender := "111111111111"
	rawJSONResp := []byte(`{
		"templates": [
			{
				"id": "111",
				"businessAccountId": 222,
				"name": "exampleName",
				"language": "en",
				"status": "APPROVED",
				"category": "ACCOUNT_UPDATE",
				"structure": {
					"header": {
						"format": "IMAGE"
					},
					"body": "example {{1}} body",
					"footer": "exampleFooter",
					"type": "MEDIA"
				}
			}
		]
	}`)
	var expectedResp models.TemplatesResponse
	err := json.Unmarshal(rawJSONResp, &expectedResp)
	require.Nil(t, err)

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, fmt.Sprintf(templatesPath, sender)))
		assert.Equal(t, fmt.Sprintf("App %s", apiKey), r.Header.Get("Authorization"))

		_, servErr := w.Write(rawJSONResp)
		assert.Nil(t, servErr)
	}))
	defer serv.Close()

	host := serv.URL
	whatsApp := whatsAppChannel{reqHandler: httpHandler{
		httpClient: http.Client{},
		baseURL:    host,
		apiKey:     apiKey,
	}}
	messageResponse, respDetails, err := whatsApp.GetTemplates(context.Background(), sender)

	require.Nil(t, err)
	assert.NotEqual(t, models.TemplatesResponse{}, messageResponse)
	assert.Equal(t, expectedResp, messageResponse)
	require.Nil(t, err)
	assert.NotNil(t, respDetails)
	assert.Equal(t, http.StatusOK, respDetails.HTTPResponse.StatusCode)
	assert.Equal(t, models.ErrorDetails{}, respDetails.ErrorResponse)
}

func TestGetTemplates4xxErrors(t *testing.T) {
	sender := "111111111111"
	tests := []struct {
		rawJSONResp []byte
		statusCode  int
	}{
		{
			rawJSONResp: []byte(`{
				"requestError": {
					"serviceException": {
						"messageId": "UNAUTHORIZED",
						"text": "Invalid login details"
					}
				}
			}`),
			statusCode: http.StatusUnauthorized,
		},
		{
			rawJSONResp: []byte(`{
				"requestError": {
					"serviceException": {
						"messageId": "TOO_MANY_REQUESTS",
						"text": "Too many requests"
					}
				}
			}`),
			statusCode: http.StatusTooManyRequests,
		},
	}
	apiKey := "secret"

	for _, tc := range tests {
		t.Run(strconv.Itoa(tc.statusCode), func(t *testing.T) {
			var expectedResp models.ErrorDetails
			err := json.Unmarshal(tc.rawJSONResp, &expectedResp)
			require.Nil(t, err)
			serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, servErr := w.Write(tc.rawJSONResp)
				assert.Nil(t, servErr)
			}))

			host := serv.URL
			whatsApp := whatsAppChannel{reqHandler: httpHandler{
				httpClient: http.Client{},
				baseURL:    host,
				apiKey:     apiKey,
			}}
			messageResponse, respDetails, err := whatsApp.GetTemplates(context.Background(), sender)
			serv.Close()

			require.Nil(t, err)
			assert.NotEqual(t, http.Response{}, respDetails.HTTPResponse)
			assert.NotEqual(t, models.ErrorDetails{}, respDetails.ErrorResponse)
			assert.Equal(t, expectedResp, respDetails.ErrorResponse)
			assert.Equal(t, tc.statusCode, respDetails.HTTPResponse.StatusCode)
			assert.Equal(t, models.TemplatesResponse{}, messageResponse)
		})
	}
}