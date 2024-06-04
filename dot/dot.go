package dot

import (
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	IGNORE = ".gitignore"
)

func DotCommand(dotPath string) error {

	dirPahts, err := GetPaths(dotPath)
	if err != nil {

		fmt.Println(err)
		return err

	}

	dirStructs := Isexe(dirPahts)
	err = CreateSymlink(dirStructs, dotPath)

	if err != nil {

		return err
	}
	return nil

}

func CreateSymlink(dotfiles []Dotfile, source string) error {

	targetPath, err := utils.GetUser("Desktop")

	if err != nil {

		return err
	}
	for _, f := range dotfiles {

		if f.IsEx {
			symlinkPath := f.Location.Name()
			sourceAbs := path.Join(source, symlinkPath)
			dest := path.Join(targetPath, symlinkPath)
			err := os.Symlink(sourceAbs, dest)

			if err != nil {
				fmt.Println("failed to create ", err)
				return err
			}

		}
	}

	return nil

}

func GetPaths(gitPath string) ([]fs.DirEntry, error) {
	dirs, err := os.ReadDir(gitPath)
	//fmt.Println(dirs)
	if err != nil {
		fmt.Println("Can't list this dir probably permissions issue ", err)
		return nil, err
	}

	toIgnore, err := readIgnore()
	if err != nil {
		fmt.Println("Can't read git ignore: ", err)
		return nil, err
	}

	var paths []fs.DirEntry
	for _, dir := range dirs {
		if !shouldIgnore(dir.Name(), toIgnore) {
			paths = append(paths, dir)
		}
	}
	//fmt.Printf("\nthih are the paths {%v}", paths)
	return paths, nil
}

func shouldIgnore(fileName string, toIgnore []string) bool {
	for _, pattern := range toIgnore {
		if match, _ := filepath.Match(pattern, fileName); match {
			//fmt.Printf("its a match {%v}  {%v} filename", pattern, fileName)
			return true

		}
	}
	return false
}

func readIgnore() ([]string, error) {

	_, err := os.Stat(IGNORE)
	if os.IsNotExist(err) {
		fmt.Println("No git ignore ")
		return nil, nil
	}

	c, err := os.ReadFile(IGNORE)
	if err != nil {
		fmt.Println("Can't read the file", err)
		return nil, err
	}

	sc := string(c)

	ignored := strings.Split(sc, "\n")
	ignored = append(ignored, ".git")
	ignored = append(ignored, IGNORE)

	//fmt.Println("this are ignored", ignored)
	return ignored, nil
}
