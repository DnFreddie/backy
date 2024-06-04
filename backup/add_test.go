package backup


import (
	"testing"
	"github.com/stretchr/testify/assert"

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
            r, err := Add_dir(&tc.item)
            assert.Nil(t, err)
            assert.Equal(t, tc.expected, r)
        })
    }
}

