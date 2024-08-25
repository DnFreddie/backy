package dot

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DnFreddie/backy/utils"
)

type Dotfile struct {
	ID         uint    `gorm:"primaryKey;autoIncrement" json:"-"`
	Location   string  `json:"location"`
	Executable bool    `json:"executable"`
	Symlink    string  `json:"symlink"`
	Absolute   string  `json:"ab_path"`
	Repo       *Repo   `gorm:"-" json:"-"`
	RepoID     *string `json:"-"`
	Ignored    bool    `json:"ignored"`
	New        bool    `json:"new"`
	Failed     string `json:"failed,omitempty"`
}

func (d *Dotfile) createSymlink(target string) error {
	d.isNew(target)
	if !d.New {
		backupDir := path.Join(d.Repo.BackupLocation, path.Base(d.Symlink))
		err := utils.Copy(d.Symlink, backupDir)
		if err != nil {
			slog.Error("Failed to rename:", "error", err)
			d.Failed = err.Error()
			return err
		}
		err = os.RemoveAll(d.Symlink)
		if err != nil {
			slog.Error("Failed to remove :", "symlink", d.Symlink, "error", err)
			return err
		}
	}

	if d.Symlink == "" {
		err := fmt.Sprintf("the symlink path shouldn't be empty")
		slog.Error("THIS SHOULDNT HAPPEN","err",err)
		d.Failed = err
		log.Fatal(err)
		return nil
	}

	err := os.Symlink(d.Absolute, d.Symlink)
	if err != nil {
		slog.Error("Failed to create symlink:", "symlink", d.Symlink, "target", d.Absolute, "error", err)
		return err
	}

	slog.Info("Symlink created successfully:", "symlink", d.Symlink, "target", d.Absolute)
	return nil
}

func (d *Dotfile) isNew(target string) {
	d.Symlink = path.Join(target, d.Location)
	if d.Symlink == "" {
		log.Fatal("Wtf is this ")
	}
	_, err := os.Stat(d.Symlink)
	if os.IsNotExist(err) {
		d.New = true
		slog.Debug("Target is new", "target", target)
	} else {
		d.New = false
		slog.Debug("Target is old", "target", target)
	}
}
func (d *Dotfile) IsExe() {
	re := regexp.MustCompile(`\.?(conf|rc)$`)
	cleanPath := re.ReplaceAllString(d.Location, "")
	d.Executable = isCmd(strings.TrimPrefix(cleanPath, "."))
}

func isCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)

	if err != nil {
		slog.Debug("Cmd is not executable", "cmd", cmd, "exec", false)
		return false
	}
	slog.Debug("Cmd is executable", "cmd", cmd, "exec", true)

	return true
}

func (d *Dotfile) ignore(toIgnore *[]string) {
	d.Ignored = false

	for _, pattern := range *toIgnore {
		if match, _ := filepath.Match(pattern, d.Location); match {
			d.Ignored = true
			slog.Debug("This should be ignonred", "name", d.Location)
			break
		}
	}
}
