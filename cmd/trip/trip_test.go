package trip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

			//err:= TripMain(TEST_PATH)
			xd, err := CreateConfig("ahshhahsahahshhashashhashh")

			fmt.Println(xd)
			assert.NoError(t, err)
		})

	}
}
