/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package trip

import (
	"errors"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"sync"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scanning the directory and checking the difference in the hash of the files",
	Long: `

Tripwire main functionality:
It scans the specified directory to detect any changes in the files' content.
The tool then outputs the current state in CSV format.

`,
	Run: func(cmd *cobra.Command, args []string) {

		err := tripScan(csvName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

var csvName string

func init() {
	scanCmd.Flags().StringVarP(&csvName, "csv", "c", "", "wirte to a given csv path")
}

func tripScan(csvPath string) error {
	done := make(chan bool)
	utils.WaitingScreen(done, "Scaning")
	confP, err := utils.Checkdir("scan_paths.json", true)
	if err != nil {
		return err
	}

	var ConfPaths []ConfigPath
	err = utils.ReadJson(confP, &ConfPaths)
	if err != nil {
		return err
	}

	if len(ConfPaths) == 0 {
		return errors.New("There are no paths in the config. First, add them with dd")
	}

	db, err := utils.InitDb(DB_PATH, nil)
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
		ScanRecursivly(fPath, db, ch)
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
		var csvName = "scan.csv"

		if csvPath != "" {
			csvName = csvPath

		}

		err = writeToCsv(&checkedArray, csvName)
		if err != nil {
			return err
		}
		done <- true
		fmt.Println("\nThe comparison scan is done, look in ", csvName)
	}

	return nil
}
