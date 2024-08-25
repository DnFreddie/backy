package dot

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/DnFreddie/backy/common"
	"github.com/DnFreddie/backy/utils"
)

func prepTest() (string, error) {
	r := &Repo{}
	force := true
	target := path.Join(os.Getenv("HOME"), "backy")
	err := r.ReadLocal(target)
	if err != nil {
		return "", fmt.Errorf("failed to read the local Repo %s: %w", r.Absolute, err)
	}

	err = r.getDots()
	if err != nil {
		return "", fmt.Errorf("error getting paths: %w", err)
	}

	slog.Debug("Started Linking")
	r.Link(force)
	slog.Debug("Saving schema")
	common.SaveStructure(r)

	return r.BackupLocation, nil
}

func TestRevert(t *testing.T) {
	const DOTS = "dotfiles"
	_, err := prepTest()
	if err != nil {
		t.Fatalf("prepTest failed: %v", err)
	}

	backupDir, err := utils.Checkdir(DOTS, false)
	r := &Repo{}
	common.Revert(r, backupDir)
	if err != nil {
		t.Fatalf("processReversion failed: %v", err)
	}

	// Check if the symlinks where correclt removed
	configDir := path.Join(os.Getenv("HOME"), ".config")
	err = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("found symlink: %s", path)
		}
		return nil
	})

	if err != nil {
		t.Errorf("error while checking for symlinks in .config: %v", err)
	}
	// Check for the existence of the LICENSE file
	licensePath := path.Join(configDir, "LICENSE")
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		t.Errorf("LICENSE file does not exist in .config")
	}

}
