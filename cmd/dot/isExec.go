package dot

import (
	"io/fs"
	"os/exec"
	"regexp"
	"strings"
)

type Dotfile struct {
	Location fs.DirEntry
	IsEx     bool
	Symlink  string
	BaseP    string
	Repo     string
	ignored  bool
}

func (d *Dotfile) IsExe() {
	re := regexp.MustCompile(`\.?(conf|rc)$`) 
	cleanPath := re.ReplaceAllString(d.Location.Name(), "") 
	d.IsEx = isCmd(strings.TrimPrefix(cleanPath,"."))
}

func isCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {

		return false
	}
	return true
}
