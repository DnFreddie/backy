package add

import (
	"fmt"
	"sync"

	"github.com/DnFreddie/backy/trip"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

func TripAdd(fPath string) error {
	db, err := utils.InitDb(trip.DB_PATH)
	if err != nil {
		return err
	}

	isNew, err := trip.CreateConfig(fPath)
	if err != nil {
		return err
	}

	if isNew {
		fmt.Printf("The %v does already exist in db. Try scan flag\n", fPath)
		return nil
	}

	var wg sync.WaitGroup
	ch := make(chan utils.FileProps)
	numWorkers := 1

	err = db.AutoMigrate(&utils.FileProps{})
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		trip.ScanRecursivly(fPath, db, ch)
		close(ch)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processBatches(ch, db)
		}()
	}

	wg.Wait()
	return nil
}

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
