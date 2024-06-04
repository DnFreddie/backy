package trip

import (
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

func processBatches(ch chan utils.FileProps, db *gorm.DB) {
	var batch []utils.FileProps
	var count int
	for item := range ch {
		count++
		batch = append(batch, item)
		if len(batch) == 100 {
			db.CreateInBatches(batch, len(batch))
			fmt.Println("Processed count during execution:", count)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		db.CreateInBatches(batch, len(batch))
		fmt.Println("Processed count during execution:", count)
	}
}
