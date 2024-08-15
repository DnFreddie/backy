package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/DnFreddie/backy/utils"
	"github.com/stretchr/testify/assert"
)

func TestDownloadRepo(t *testing.T) {
	testCases := []struct {
		name    string
		zipUrl  string
		repPath string
		err     bool
	}{
		//Don't change the order else clean up the created file if it exists
		{"Wrong URL", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/wrong.zip", "DnFreddie", true},
		{"Wrong Path", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/main.zip", "/test/xd/name", true},
		{"Download correct archive", "https://github.com/DnFreddie/DnFreddie/archive/refs/heads/main.zip", "DnFreddie", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := Repo{
				zipUrl:   tc.zipUrl,
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
func TestLink(t *testing.T) {
	testCases := []struct {
		name     string
		err      bool
		r        Repo
		symlinks []string
	}{
		{"Test Link", false, Repo{}, []string{"utils", "cmd", "LICENSE"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			err := tc.r.GetInfo("/home/test/backy")
			if (err != nil) != tc.err {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tc.err)
				return
			}

			err = tc.r.getDots()
			if err != nil {
				t.Errorf("getDots() error = %v", err)
				return
			}

			tc.r.createBackup()
			tc.r.Link(true)

			for _, i := range tc.symlinks {
				symlinkPath, err := utils.GetUser(filepath.Join(".config", i))
				assert.NoError(t, err)


				assertSymlink, err := os.Lstat(symlinkPath)
				assert.NoError(t, err)

				if assertSymlink == nil {
				t.Fatal("assertSymlink is nil")
}
				if assertSymlink.Mode()&os.ModeSymlink != 0 {
					originFile, err := os.Readlink(symlinkPath)
					assert.NoError(t, err)

					fmt.Println("Resolved symlink to: ", originFile)
				} else {
					t.Errorf("Expected symlink at %s, but it is not a symlink", symlinkPath)
				}
			}
		})
	}
}

