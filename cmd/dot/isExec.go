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
}
// Without bash -c it won't work on nixos

func Isexe(dirs []fs.DirEntry) []Dotfile {
	var list []Dotfile
	for _, d := range dirs {
		re := regexp.MustCompile(`(rc|\.conf)$`)
		cleanPath := re.ReplaceAllString(strings.TrimPrefix(d.Name(), "."), "")

		dotfile := Dotfile{
			Location: d,

			IsEx: isCommandAvailable(cleanPath),
		}

		list = append(list, dotfile)

	}

	return list

}

func isCommandAvailable(name string) bool {

	cmd := exec.Command("bash", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
