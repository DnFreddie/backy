package dot

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Dotfile struct {
	ID         uint     `gorm:"primaryKey;autoIncrement" json:"-"`
	Location   string   `json:"location"`
	Executable bool     `json:"executable"`
	Symlink    string   `json:"symlink"`
	AbPath     string   `json:"ab_path"`
	Repo       *Repo    `gorm:"-" json:"-"` 
	RepoID     *string  `json:"-"`        
	Ignored    bool     `json:"ignored"`
	New        bool     `json:"new"`
	Failed     *string  `json:"failed,omitempty"`         
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
