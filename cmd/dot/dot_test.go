package dot

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUrl(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected bool
	}{ {"wrong url", "www.xdxdasdasdasdasd", false},
		{"correct url", "https://github.com/DnFreddie/Notes", true},
		{"ssh url", "git@github.com:DnFreddie/Notes.git", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, isUrl(tc.url))
		})
	}
}

func TestGetHeadUrl(t *testing.T) {
	testCases := []struct {
		name   string
		url    string
		err    bool
		zipUrl string
	}{
		{"wrong url", "www.xdxdasdasdasdasd", true, ""},
		{"correct url", "https://github.com/DnFreddie/Notes", false, "https://github.com/DnFreddie/Notes/archive/refs/heads/hugo.zip"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zipUrl, err := getHeadUrl(tc.url)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.zipUrl, zipUrl)
			}
		})
	}
}


func TestIsExe(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		isExe bool
	}{
		{"Regular command", "ls", true},
		{"Dotfile with .conf", "ls.conf", true},
		{"Dotfile with rc", ".ls.rc", true},
		{"Non-command", "testR", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := &customDirEntry{name: tt.dirName}
			dot := &Dotfile{Location: location}

			dot.IsExe()

			assert.Equal(t, tt.isExe, dot.IsEx)
		})
	}
}
type customDirEntry struct {
	name string
}

func (c *customDirEntry) Name() string {
	return c.name
}

func (c *customDirEntry) Type() os.FileMode {
	return 0 
}

func (c *customDirEntry) Info() (os.FileInfo, error) {
	return nil, nil 
}

func (c *customDirEntry) IsDir() bool {
	return false 
}
