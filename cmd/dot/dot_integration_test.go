package dot

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	//"github.com/stretchr/testify/assert"
)

func TestDownloadRepo(t *testing.T) {
	testCases := []struct {
		name    string
		zipUrl  string
		repPath string
		err     bool
	}{
		{"Wrong URL", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/wrong.zip", "DnFreddie", true},
		{"Wrong Path", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/main.zip", "/test/xd/name", true},
		{"Download correct archive", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/main.zip", "DnFreddie", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := Repo{
				zipUrl:  tc.zipUrl,
				RepoName: tc.repPath,
			}

			err := repo.downloadRepo()

			if err != nil {
				assert.True(t, tc.err, "this should error")
				zipFile := tc.repPath + ".zip"
				_, err = os.Stat(zipFile)
				assert.True(t, os.IsNotExist(err), "File should not exist")
			} else {
				assert.False(t, tc.err, "This should be successful")
				zipFile := tc.repPath + ".zip"
				_, err = os.Stat(zipFile)
				assert.NoError(t, err, "File should exist")
			}

			// Clean up the created file if it exists
		})
	}
}
var cRepos = Repo{
	RepoName:      "roles",
	DefaultBranch: "main",
}

func TestGitClone(t *testing.T) {
	testCases := []struct {
		name       string
		url        string
		err        bool
		r          Repo
		resultRepo Repo
	}{
		{"Wrong URL", "wrong_url", true, Repo{}, cRepos},
		{"Valid URL", "https://github.com/DnFreddie/roles", false, Repo{}, cRepos},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("Running test:", tc.name)
			tc.r.Url = tc.url
			err := tc.r.Clone(tc.r.Url)
			if err != nil {
				assert.Equal(t, tc.err, true)
			} else {
				assert.Equal(t, tc.r.Url, tc.url)
				assert.Equal(t, tc.r.RepoName, cRepos.RepoName)
				assert.Equal(t, tc.r.DefaultBranch, cRepos.DefaultBranch)
			}
		})
	}
}
