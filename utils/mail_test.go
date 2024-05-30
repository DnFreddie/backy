package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadAndSendEmail(t *testing.T) {
	testCases := []struct {
		name     string
		expected Email_Creds
		body     string
		wantErr  bool
	}{
		{
			name: "Correct credentials and message",
			expected: Email_Creds{
				Email:  "szopen_test@gmail.com",
				Passwd: "12344",
			},
			body:    "This is a test email message",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			 err := tc.expected.readTheConfig()
			if assert.NoError(t, err, "Error reading the config file") {
				assert.Equal(t, tc.expected.Email, tc.expected.Email)
				assert.Equal(t, tc.expected.Passwd, tc.expected.Passwd)
				err = SendMessage(tc.body,	&tc.expected )
				if tc.wantErr {
					assert.Error(t, err, "Expected an error while sending message")
				} else {
					assert.NoError(t, err, "Expected no error while sending message")
				}
			}
		})
	}
}
