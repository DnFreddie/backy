package backup

import (
	"fmt"
	"path"

	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

func addDir(paths *[]string) ([]string, error) {
	var newPaths []string

	for _, p := range *paths {
		new_path, err := utils.MakeAbsoulute(p)
		if err != nil {
			fmt.Println(p, "Doesn't exist")

			continue
		}

		newPaths = append(newPaths, new_path)
	}

	return newPaths, nil
}

const (
	BACK_PATH = "backy_back.sql"
)

type Brecord struct {
	*gorm.Model
	TargetPath string
	CurrPath   string
}

func addPaths(FDirs []string) error {

	db, err := utils.InitDb(BACK_PATH, &Brecord{})

	if err != nil {
		return err
	}

	for _, f := range FDirs {

		record := Brecord{
			TargetPath: f,
			CurrPath:   f,
		}
		fmt.Println(path.Base(f), "was successfully added")

		db.Create(&record)
	}

	return nil

}
