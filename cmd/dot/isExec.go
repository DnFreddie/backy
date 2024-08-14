package dot

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Dotfile struct {
	ID         uint     `gorm:"primaryKey;autoIncrement"`
	Location   string
	Executable bool
	Symlink    string
	AbPath     string
	Repo       *Repo  `gorm:"-"`
	RepoID     string  
	Ignored    bool
	New        bool
	Failed  *string
	}


func (d *Dotfile) IsExe() {
	re := regexp.MustCompile(`\.?(conf|rc)$`)
	cleanPath := re.ReplaceAllString(d.Location, "")
	d.Executable = isCmd(strings.TrimPrefix(cleanPath, "."))
}

func isCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {

		return false
	}
	return true
}

func (d *Dotfile) ignore(toIgnore *[]string) {
	d.Ignored = false

	for _, pattern := range *toIgnore {
		if match, _ := filepath.Match(pattern, d.Location); match {
			d.Ignored = true
			break
		}
	}
}
