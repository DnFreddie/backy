package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPath(t *testing.T) {
	testCases := []struct {
		name     string
		item     []string
		expected []string
	}{
		{
			name:     "Correct paths array",
			item:     []string{"wrong", "Documents", "/etc/"},
			expected: []string{"/etc"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := addDir(&tc.item)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, r)
		})
	}
}
