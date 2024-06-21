package deamon

import (
	"fmt"
	"github.com/DnFreddie/backy/backup"
	"github.com/DnFreddie/backy/utils"
	"log"
	"time"
)

func loopEnternely() {
	var record []backup.Brecord

	db, err := utils.InitDb(backup.BACK_PATH, backup.Brecord{})

	if err != nil {
		log.Fatal("Failed to connect to the db")
	}
	for {
		db.Find(&record)

		for _, i := range record {
			fmt.Println(i.TargetPath)

		}

		time.Sleep(20 * time.Second)
		fmt.Println("Looppign")

	}

}
