package scan

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/DnFreddie/backy/trip"
	"github.com/DnFreddie/backy/utils"
)

func TripScan(csvPath string) error {
	confP, err := utils.Checkdir("scan_paths.json", true)
	if err != nil {
		return err
	}

	var ConfPaths []trip.ConfigPath
	err = utils.ReadJson(confP, &ConfPaths)
	if err != nil {
		return err
	}

	if len(ConfPaths) == 0 {
		return errors.New("There are no paths in the config. First, add them with TripAdd")
	}

	db, err := utils.InitDb(trip.DB_PATH)
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
		trip.ScanRecursivly(fPath, db, ch)
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
