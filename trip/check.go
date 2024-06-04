package trip

import (
	"bytes"
	"fmt"
	"reflect"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

type Status string

const (
    NotExist Status = "Not Exists"
    Changed Status = "Changed"
)

type Compared struct {
    Status
    FilePath string
    Dir      string
}
func checkSumFiles(ch chan utils.FileProps,  db *gorm.DB, fchan chan Compared) {
	for item := range ch {
		user:= utils.FileProps{}
		 comperdItem:= Compared{
			FilePath: item.FilePath,
			Dir: item.DirPath,
		}

		db.Where("file_path = ?", item.FilePath).Find(&user)

		if reflect.ValueOf(user).IsZero() {
			comperdItem.Status=NotExist
			fchan <- comperdItem
		} else {
			compersion := compere(user, item)
			if compersion != "" {
				comperdItem.Status=Changed
				fchan <- comperdItem
			}

		}
	}

}

func compere(saved utils.FileProps, newI utils.FileProps) string {

	if !bytes.Equal(saved.Hash, newI.Hash) {
		fmt.Println("the hasesh doens't work ", saved.FilePath)
		return saved.FilePath
	}

	return ""

}
