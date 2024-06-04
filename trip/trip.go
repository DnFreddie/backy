package trip

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"sync"
)

func TripAdd(fPath string) error {
	db, err := utils.InitDb()
	if err != nil {
		return err
	}

	isNew, err := createConfig(fPath)
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
		scanRecursivly(fPath, db, ch)
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

func TripScan(csvPath string) error {
	confP, err := utils.Checkdir("scan_paths.json")
	if err != nil {
		return err
	}

	var ConfPaths []ConfigPath
	err = utils.ReadJson(confP, &ConfPaths)
	if err != nil {
		return err
	}

	if len(ConfPaths) == 0 {
		return errors.New("There are no paths in the config. First, add them with TripAdd")
	}

	db, err := utils.InitDb()
	if err != nil {
		return err
	}

	checked := make(chan Compared)
	ch := make(chan utils.FileProps)
	var checkedArray []Compared
	fPath := ConfPaths[0].Fpath
	var wg sync.WaitGroup
	numWorkers := 5

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanRecursivly(fPath, db, ch)
		close(ch)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			checkSumFiles(ch, db, checked)
		}()
	}

	go func() {
		wg.Wait()
		close(checked)
	}()

	for i := range checked {
		checkedArray = append(checkedArray, i)
	}

	if len(checkedArray) == 0 {
		slog.Info("Everything is fine with this dir")
	} else {
		var csvName = "trip_scan.csv"

		if csvPath != "" {
			csvName = csvPath

		}

		err = writeToCsv(&checkedArray, csvName)
		fmt.Println("The comparison scan is done, look in ", csvName)
		if err != nil {
			return err
		}
	}

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
			if err := writer.Write(formatData); err != nil {
				return fmt.Errorf("failed to write data: %w", err)
			}
		}
	}
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
