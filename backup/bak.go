package backup

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/DnFreddie/backy/utils"
)

func Back(pathsArr *[]string) error {
	nowT := time.Now().Format("20060102150405")
	dirPath := path.Join("backups", nowT)
	dest, err := utils.Checkdir(dirPath, false)
	if err != nil {
		return err
	}
	paths, err := addDir(pathsArr)

	for _, p := range paths {
		stat, err := os.Stat(p)
		if err != nil {
			fmt.Println("Error stating file:", err)
			return err
		}

		finalDest := path.Join(dest, path.Base(p))
		if stat.IsDir() {
			err = copyDir(p, finalDest)

		} else {
			err = copyFile(p, finalDest)
		}

		err = writeSchema(&paths, dest)
		if err != nil {
			os.RemoveAll(dirPath)
			fmt.Println("Something went wrong:", err)
			return err
		}
		fmt.Println(path.Base(p), "has been backed")
	}

	return nil

}
func writeSchema(bPaths *[]string, dirPath string) error {
	data := strings.Join(*bPaths, "\n")
	dest := path.Join(dirPath, "schema.txt")

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}
