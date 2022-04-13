package email

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/infobip-community/infobip-api-go-sdk/internal"
	"github.com/infobip-community/infobip-api-go-sdk/pkg/infobip/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetDeliveryReports(t *testing.T) {
	apiKey := "secret"

	rawJSONResp := []byte(`
		{
			"results": [
				{
					 "bulkId": "test-bulk-583",
					 "messageId": "test-message-355",
					 "to": "joan.doe0@example.com",
					 "sentAt": "2021-11-08T21:49:44.772+0000",
					 "doneAt": "2021-11-08T21:49:49.930+0000",
					 "messageCount": 1,
					 "price": {
						  "pricePerMessage": 0,
						  "currency": "EUR"
					 },
					 "status": {
						  "groupId": 3,
						  "groupName": "DELIVERED",
						  "id": 5,
						  "name": "DELIVERED_TO_HANDSET",
						  "description": "Message delivered to handset"
					 },
					 "error": {
						  "groupId": 0,
						  "groupName": "OK",
						  "id": 0,
						  "name": "NO_ERROR",
						  "description": "No Error",
						  "permanent": false
					 },
					 "channel": "EMAIL"
				},
				{
					 "bulkId": "test-bulk-748",
					 "messageId": "test-message-411",
					 "to": "joan.doe0@example.com",
					 "sentAt": "2021-11-08T21:49:43.854+0000",
					 "doneAt": "2021-11-08T21:49:50.734+0000",
					 "messageCount": 1,
					 "price": {
						  "pricePerMessage": 0,
						  "currency": "EUR"
					 },
					 "status": {
						  "groupId": 3,
						  "groupName": "DELIVERED",
						  "id": 5,
						  "name": "DELIVERED_TO_HANDSET",
						  "description": "Message delivered to handset"
					 },
					 "error": {
						  "groupId": 0,
						  "groupName": "OK",
						  "id": 0,
						  "name": "NO_ERROR",
						  "description": "No Error",
						  "permanent": false
					 },
					 "channel": "EMAIL"
				}
			]
		}
	`)

	var expectedResp models.EmailDeliveryReportsResponse
	err := json.Unmarshal(rawJSONResp, &expectedResp)
	require.NoError(t, err)

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, getDeliveryReportsPath))
		assert.Equal(t, fmt.Sprint("App ", apiKey), r.Header.Get("Authorization"))

		_, servErr := w.Write(rawJSONResp)
		assert.Nil(t, servErr)
	}))
	defer serv.Close()

	email := Channel{ReqHandler: internal.HTTPHandler{
		HTTPClient: http.Client{},
		BaseURL:    serv.URL,
		APIKey:     apiKey,
	}}

	queryParams := make(map[string]string)

	resp, respDetails, err := email.GetDeliveryReports(context.Background(), queryParams)

	require.NoError(t, err)
	assert.NotEqual(t, models.EmailDeliveryReportsResponse{}, resp)
	assert.Equal(t, expectedResp, resp)
	assert.NotNil(t, respDetails)
	assert.Equal(t, http.StatusOK, respDetails.HTTPResponse.StatusCode)
	assert.Equal(t, models.ErrorDetails{}, respDetails.ErrorResponse)
}
