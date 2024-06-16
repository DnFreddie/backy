package dot

import (
	"encoding/csv"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const (
	IGNORE     = ".gitignore"
	BACK_CONF  = "back_conf"
	REVERT_CSV = "schema.csv"
	TARGET     = "Desktop"
)

func DotCommand(repo string) error {
	var URL bool
	var dest string

	if strings.Contains(repo, "git@") {
		URL = true
	} else {
		URL = isUrl(repo)

	}
	if URL {
		clonedDest, err := gitClone(repo)
		if err != nil {
			log.Fatal("Failed to copy url")
		}
		dest = clonedDest
	} else {
		dest = repo
	}

	if !path.IsAbs(dest) {

		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Therse smth wrong with this directroy check perrmisons ")
		}
		dest = path.Join(pwd, dest)
	}

	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		log.Fatalf("%v doesn't exist\n", path.Base(dest))
	}

	dirPaths, err := GetPaths(dest)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dirStructs := Isexe(dirPaths)

	err = createSymlink(dirStructs, dest)
	if err != nil {
		return err
	}

	return nil
}

func createTempBack(source string, backupDir string, csvF *csv.Writer, sourceAbs string, newDest string) (bool, error) {

	_, err := os.Stat(source)

	if os.IsNotExist(err) {
		return false, nil
	}

	dest := path.Join(backupDir, path.Base(source))
	fmt.Println("Already exist", path.Base(dest))
	err = os.Rename(source, dest)

	if err != nil {
		return false, err
	}
	data := [][]string{
		{source, dest},
	}
	err = csvF.WriteAll(data)

	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	err = os.Symlink(sourceAbs, newDest)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, nil

}

func createSymlink(dotfiles []Dotfile, source string) error {

	targetPath, err := utils.GetUser(TARGET)
	if err != nil {

		return err
	}

	nowT := time.Now().Format("20060102150405")

	backupDir, err := utils.Checkdir(path.Join(BACK_CONF, nowT), false)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(path.Join(backupDir, REVERT_CSV))
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(f)

	for _, f := range dotfiles {

		if f.IsEx {
			symlinkPath := f.Location.Name()
			sourceAbs := path.Join(source, symlinkPath)
			dest := path.Join(targetPath, symlinkPath)

			wasCreated, err := createTempBack(dest, backupDir, writer, sourceAbs, dest)
			if err != nil {
				fmt.Println("Error creating temporary backup:", err)
				return err
			}

			if !wasCreated {

				err = os.Symlink(sourceAbs, dest)

				if err != nil {
					fmt.Println("Failed to create symlink:", err)
					return err

				}
				fmt.Println("Created symilnk", path.Base(dest))

				data := [][]string{
					{dest, "new"},
				}

				err = writer.WriteAll(data)

				if err != nil {
					return err
				}

			}
		}
	}

	return nil

}

func GetPaths(gitPath string) ([]fs.DirEntry, error) {
	dirs, err := os.ReadDir(gitPath)
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
	return paths, nil
}
