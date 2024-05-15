package dot

import (
	"io/fs"
	"os/exec"
)
// Without bash -c it won't work on nixos
func isCommandAvailable(name string) bool {
	cmd := exec.Command("bash", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}




func Isexe(dirs []fs.DirEntry) []Dotfile  {

	var list []Dotfile
	for _,d := range dirs{
	dotfile :=  Dotfile{
			Location: d ,
			IsEx: isCommandAvailable(d.Name()),
			

		}

	list = append(list, dotfile)

	}

	return list

}
