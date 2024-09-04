package backup

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/DnFreddie/backy/hash"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
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
	fmt.Println("This are the scanned path" ,paths)
	err = scanPaths(paths)
	if err != nil {
		return err
	}

	return nil
}

func addDir(paths *[]string) ([]string, error) {
	var newPaths []string

	for _, p := range *paths {
		new_path, err := utils.MakeAbsolute(p)
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

type Brecord_old struct {
	*gorm.Model
	TargetPath string `gorm:"unique"`
	CurrPath   string
}

func scanPaths(absPaths []string) error {
	for _, i := range absPaths {
		node, err := hash.WalkDir(i)
		if err != nil {
			slog.Error("Node errored", "err", err)
			fmt.Println(err)
			continue 
		}

		backDir,err := utils.Checkdir(BACKUP_DIR,false)

		if err != nil{
				log.Fatal(err)
		}
		fmt.Println("this is the back dir",backDir)

		if err != nil{
			return err
		}
		backup := path.Join(backDir,path.Base(i) + ".json" )
		f, err := os.Create(backup)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		
		defer f.Close()

		jsn, err := json.Marshal(node)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}

		_, err = f.Write(jsn)
		if err != nil {
			log.Fatalf("Failed to write JSON to file: %v", err)
		}
	}

	return nil
}

