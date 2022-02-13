package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidStickerMessage(t *testing.T) {
	tests := []struct {
		name     string
		instance StickerMessage
	}{
		{
			name: "minimum input",
			instance: StickerMessage{
				MessageCommon: MessageCommon{
					From: "16175551213",
					To:   "16175551212",
				},
				Content: StickerContent{MediaURL: "https://www.mypath.com/audio.mp3"},
			}},
		{
			name: "complete input",
			instance: StickerMessage{
				MessageCommon: GenerateTestMessageCommon(),
				Content: StickerContent{
					MediaURL: "https://www.mypath.com/audio.mp3",
				},
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

func TestStickerMessageConstraints(t *testing.T) {
	msgCommon := GenerateTestMessageCommon()
	tests := []struct {
		name    string
		content StickerContent
	}{
		{
			name:    "missing Content MediaURL",
			content: StickerContent{},
		},
		{
			name:    "Content MediaURL too long",
			content: StickerContent{MediaURL: fmt.Sprintf("https://www.g%sgle.com", strings.Repeat("o", 2040))},
		},
		{
			name:    "Content invalid MediaURL",
			content: StickerContent{MediaURL: "asd"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := StickerMessage{
				MessageCommon: msgCommon,
				Content:       tc.content,
			}
			err := msg.Validate()
			require.NotNil(t, err)
		})
	}
}