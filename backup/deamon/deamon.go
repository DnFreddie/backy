package deamon

import (
	"fmt"
	"github.com/DnFreddie/backy/backup"
	"github.com/DnFreddie/backy/utils"
	"github.com/robfig/cron/v3"
	"log"
)

func loopEternally(croneRecord string) {
	db, err := utils.InitDb(backup.BACK_PATH, backup.Brecord{})
	if err != nil {
		log.Fatal("Failed to connect to the db:", err)
	}

	c := cron.New()

	_, err = c.AddFunc(croneRecord, func() {
		var records []backup.Brecord
		var paths []string

		
		result := db.Find(&records)
		if result.Error != nil {
			log.Println("Failed to retrieve records:", result.Error)
			return
		}

		for _, record := range records {
			paths = append(paths, record.TargetPath)

		}
		fmt.Println(paths)
		err:= backup.Back(&paths)
		if  err != nil{
			fmt.Println(err)
			return
		}

	})
	if err != nil {
		log.Fatal("Failed to add cron job:", err)
	}

	c.Start()

	select {}
}
