package backup

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"time"
	"github.com/DnFreddie/backy/utils"
)

func Back(pathsArr *[]string,backup_path string) error {
	nowT := time.Now().Format("20060102150405")
	dirPath := path.Join(backup_path, nowT)
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
	var records [][]string

	schema := path.Join(dirPath, utils.SCHEMA_CSV)

	f, err := os.Create(schema)
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	for _, p := range *bPaths {
		backP := path.Join(dirPath, path.Base(p))
		records = append(records, []string{p, backP})
	}

	err = csvWriter.WriteAll(records)
	if err != nil {
		return err
	}


	return nil
}
