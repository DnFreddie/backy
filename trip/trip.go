package trip

import (
	"encoding/csv"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"sync"
)

func TripMain(fPath string) error {
	var checkedArray []Compared
	ch := make(chan utils.FileProps)
	checked := make(chan Compared)
	var wg sync.WaitGroup
	db, err := utils.InitDb()
	if err != nil {
		return err
	}

	numWorkers := 20
	wg.Add(numWorkers)

	if !db.Migrator().HasTable(&utils.FileProps{}) {
		err = db.AutoMigrate(&utils.FileProps{})

		if err != nil {
			return err
		}

		go func() {
			scanRecursivly(fPath, db, ch)
			close(ch)
		}()

		for i := 0; i < numWorkers; i++ {
			go processBatches(ch, &wg, db)
		}

		wg.Wait()
	} else {
		go func() {
			scanRecursivly(fPath, db, ch)
			close(ch)
		}()

		for i := 0; i < numWorkers; i++ {
			go checkSumFiles(ch, &wg, db, checked)
		}

		go func() {
			wg.Wait()
			close(checked)
		}()
		for i := range checked {
			checkedArray = append(checkedArray, i)

		}

		wg.Wait()
		err = writeToCsv(&checkedArray, "test.csv")
	}
	if err != nil {
		fmt.Println("smth went wrong wiht the csv")
		return err

	}

	var count int64
	if err := db.Model(&utils.FileProps{}).Count(&count).Error; err != nil {
		return err
	}
	fmt.Println("Total count of FileProps:", count)

	return nil
}

func writeToCsv(data *[]Compared, filePath string) error {
	csvFile := filePath
	_, err := os.Stat(csvFile)
	if os.IsNotExist(err) {
		f, err := os.Create(csvFile)
		defer f.Close()

		if err != nil {
			return fmt.Errorf("can't create the file %v due to: %v", csvFile, err)
		}

		writer := csv.NewWriter(f)
		data := [][]string{
			{"status", "directory", "file_path"},
		}
		err = writer.WriteAll(data)

		if err != nil {
			slog.Error("Can't the record for the file ", csvFile, err)
			return err
		}
	} else {

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		writer := csv.NewWriter(f)
		defer writer.Flush()

		for _, item := range *data {
			formatData := []string{string(item.Status), item.Dir, item.FilePath}
			fmt.Println(formatData)
			if err := writer.Write(formatData); err != nil {
				return fmt.Errorf("failed to write data: %w", err)
			}
		}

	}
	return nil
}

func processBatches(ch chan utils.FileProps, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	var batch []utils.FileProps
	var count int

	for item := range ch {
		count++
		batch = append(batch, item)
		if len(batch) == 300 {
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
