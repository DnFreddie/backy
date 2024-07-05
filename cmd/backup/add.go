package backup

import (
	"fmt"
	"path"

	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)



func Add_command(args *[]string) error {
	paths, err := addDir(args)

	if err != nil {
		return err
	}
	if len(paths)==0{
		fmt.Println("Skipping... Nothing to add ")
		return nil
	}
	err = addPaths(paths)
	if err != nil {
		return err
	}

	return nil
}

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
	BACKUP_DB = "backy_back.sql"
)

type Brecord struct {
	*gorm.Model
	TargetPath string `gorm:"unique"`
	CurrPath   string
}

// https://gorm.io/docs/create.html#Upsert-x2F-On-Conflict
func addPaths(FDirs []string) error {
	db, err := utils.InitDb(BACKUP_DB, &Brecord{})
	if err != nil {
		return err
	}

	for _, f := range FDirs {
		record := Brecord{
			TargetPath: f,
			CurrPath:   path.Base(f),
		}
		fmt.Println(path.Base(f), "was successfully added")

		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "target_path"}}, 
			DoUpdates: clause.Assignments(map[string]interface{}{"target_path": f}), 
		}).Create(&record)
	}

	return nil
}

