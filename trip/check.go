package trip

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

func checkSumFiles(ch chan utils.FileProps, wg *sync.WaitGroup, db *gorm.DB, fchan chan string) {
	defer wg.Done()
	for item := range ch {
		var user utils.FileProps
		db.Where("file_path = ?", item.FilePath).Find(&user)
		if reflect.ValueOf(user).IsZero() {
			fchan <- item.FilePath
		} else {
			compersion := compere(user, item)
			if compersion != "" {
				fchan <- fmt.Sprintf("comeperion is this %v", compersion)
			}
			user = utils.FileProps{}

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
