package dot

import (
	"io/fs"
	"os/exec"
	"regexp"
	"strings"
)

// Without bash -c it won't work on nixos
func isCommandAvailable(name string) bool {

	cmd := exec.Command("bash", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

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
