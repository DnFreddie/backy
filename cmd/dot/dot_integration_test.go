package dot

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadGit(t *testing.T) {
	testCases := []struct {
		name     string
		zipUrl   string
		repPath  string
		expected bool
	}{
		{"Download correct archive", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/main.zip", "DnFreddie", false},
		{"Wrong URL", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/wrong.zip", "DnFreddie", true},
		{"Wrong Path", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/main.zip", "/test/xd/name", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := downloadGit(tc.zipUrl, tc.repPath)

			if (err != nil) != tc.expected {
				if tc.expected {
					t.Errorf("expected an error, but got none")
				} else {
					t.Errorf("expected no error, but got: %v", err)
				}
			}

			if err == nil {
				if _, err := os.Stat(tc.repPath); os.IsNotExist(err) {
					t.Errorf("expected directory %s to exist, but it does not", tc.repPath)
				}
			} else {
				if _, err := os.Stat(tc.repPath + "zip"); os.IsNotExist(err) {
					assert.NotErrorIs(t, err, nil)
				}
			}
		})
	}
}
