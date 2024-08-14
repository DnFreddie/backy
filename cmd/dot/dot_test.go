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
		name string
		url string
		err bool
		archvieUrl string
		r Repo
		corect Repo


	}{
		{"wrong url", "www.xdxdasdasdasdasd", true, "",Repo{} ,cRepos},
		{"correct url", "https://github.com/DnFreddie/Notes", false, "https://github.com/DnFreddie/Notes/archive/refs/heads/hugo.zip",Repo{},cRepos},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			 err := tc.r.getHeadUrl(tc.url)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.r.zipUrl, tc.archvieUrl)
			}
		})
	}
}


func TestIsExe(t *testing.T) {
	tests := []struct {
		name     string
		fileName  string
		isExe bool
	}{
		{"Regular command", "ls", true},
		{"Dotfile with .conf", "ls.conf", true},
		{"Dotfile with rc", ".ls.rc", true},
		{"Non-command", "testR", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dot := &Dotfile{Location: tc.fileName}

			dot.IsExe()

			assert.Equal(t, tc.isExe, dot.Executable)
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
