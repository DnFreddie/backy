package dot

import (
	"fmt"
	"github.com/DnFreddie/backy/common"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"testing"
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
	bLocation, err := prepTest()
	if err != nil {
		t.Fatalf("prepTest failed: %v", err)
	}

	baseDir := filepath.Dir(bLocation)
	newFileName := bLocation[len(baseDir)+1:]

	newPath := filepath.Join(baseDir, newFileName)
	r, err := processReversion(newPath)
	if err != nil {
		t.Fatalf("processReversion failed: %v", err)
	}

	for _, d := range *r.Dots {
		if err := d.cleanLinks(); err != nil {
			fmt.Println(bLocation)
			t.Errorf("failed to clean links: %v", err)
			continue
		}

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

	err = os.RemoveAll(bLocation)
	if err != nil {
		t.Errorf("failed to remove path %s: %v", bLocation, err)
	}
}
