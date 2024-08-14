package utils
import (
	"github.com/stretchr/testify/assert"
	"testing"
)


type email_creds_test struct{
Email string
Passwd string
}
func TestReadAndSendEmail(t *testing.T) {
	testCases := []struct {
		name     string
		expected   email_creds_test
		body     string
		err      bool
	}{
		{
			name: "Wrong credentials and message",
			expected: email_creds_test{
				Email:  "szopen_test@gmail",
				Passwd: "12344",
			},
			body:    "This is a test email message",
			err:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SendMessage(tc.body, tc.expected.Email, tc.expected.Passwd)
			if err != nil {
				assert.Equal(t, tc.err, true)
			} else {
				assert.Equal(t, tc.err, false)
			}
		})
	}
}

