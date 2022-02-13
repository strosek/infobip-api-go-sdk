package models

import (
	"infobip-go-client/pkg/infobip/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidLocationMessage(t *testing.T) {
	msgCommon := GenerateTestMessageCommon()
	tests := []struct {
		name     string
		instance LocationMessage
	}{
		{
			name: "minimum input",
			instance: LocationMessage{
				MessageCommon: MessageCommon{
					From: "16175551213",
					To:   "16175551212",
				},
				Content: LocationContent{Latitude: utils.Float32Ptr(10.5), Longitude: utils.Float32Ptr(10.5)},
			},
		},
		{
			name: "complete input",
			instance: LocationMessage{
				MessageCommon: msgCommon,
				Content: LocationContent{
					Name:      "Some Place",
					Address:   "My Address",
					Latitude:  utils.Float32Ptr(10.5),
					Longitude: utils.Float32Ptr(10.5),
				},
			},
		},
		{
			name: "Latitude and longitude 0",
			instance: LocationMessage{
				MessageCommon: msgCommon,
				Content:       LocationContent{Latitude: utils.Float32Ptr(0), Longitude: utils.Float32Ptr(0)},
			},
		},
		{
			name: "Latitude and longitude edge values positive",
			instance: LocationMessage{
				MessageCommon: msgCommon,
				Content:       LocationContent{Latitude: utils.Float32Ptr(90), Longitude: utils.Float32Ptr(180)},
			},
		},
		{
			name: "Latitude and longitude edge values negative",
			instance: LocationMessage{
				MessageCommon: msgCommon,
				Content:       LocationContent{Latitude: utils.Float32Ptr(-90), Longitude: utils.Float32Ptr(-180)},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.instance.Validate()
			require.Nil(t, err)
		})
	}
}

func TestTextLocationConstraints(t *testing.T) {
	msgCommon := GenerateTestMessageCommon()
	tests := []struct {
		name    string
		content LocationContent
	}{
		{
			name:    "empty Content field",
			content: LocationContent{},
		},
		{
			name:    "missing Content Latitude",
			content: LocationContent{Longitude: utils.Float32Ptr(10.5)},
		},
		{
			name:    "invalid Content Latitude",
			content: LocationContent{Latitude: utils.Float32Ptr(91), Longitude: utils.Float32Ptr(10.5)},
		},
		{
			name:    "missing Content Longitude",
			content: LocationContent{Latitude: utils.Float32Ptr(10.5)},
		},
		{
			name:    "invalid Content Longitude",
			content: LocationContent{Latitude: utils.Float32Ptr(10.5), Longitude: utils.Float32Ptr(181)},
		},
		{
			name:    "Content Name too long",
			content: LocationContent{Name: strings.Repeat("a", 1001)},
		},
		{
			name:    "Content Address too long",
			content: LocationContent{Address: strings.Repeat("a", 1001)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := LocationMessage{
				MessageCommon: msgCommon,
				Content:       tc.content,
			}
			err := msg.Validate()
			require.NotNil(t, err)
		})
	}
}