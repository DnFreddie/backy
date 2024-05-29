package trip

import (
	"testing"
)

const (
	TEST_PATH = "/home/aura/.config/"
)

func TestAddPath(t *testing.T) {
	testCases := []struct {
		name     string
		item     []string
		expected []string
	}{
		{
			name:     "Correct paths array",
			item:     []string{"xd", "Documents", "/etc/"},
			expected: []string{"/etc"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test,err:= TripCommand(TEST_PATH)
			
			if err != nil {
				return  
			}
			printDirectory(test)
		})

	}
}
