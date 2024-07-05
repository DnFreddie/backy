package revert

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func revertBackups(backupDir string) error {
	confDir, err := utils.Checkdir(backupDir, false)
	if err != nil {
		return err
	}
	options, err := os.ReadDir(confDir)
	if err != nil {
		return err
	}

	if len(options) < 1 {
		fmt.Println("No backups to revert")
		return nil
	}

	revDir, err := chooseBackupVersion(options)
	if err != nil {
		return err
	}

	chosenPath := path.Join(confDir, revDir.Name())
	if err := processReversion(chosenPath); err != nil {
		fmt.Println(err)
		return err
	}

	err = os.RemoveAll(chosenPath)
	if err != nil {
		log.Fatal("Have you changed permission? This shouldn't have happened")
	}
	fmt.Println("succesfully removed ", chosenPath)

	fmt.Println("The backup reversion was successful")
	return nil
}

func processReversion(chosenPath string) error {

	csvPath := path.Join(chosenPath, utils.SCHEMA_CSV)
	_, err := os.Stat(csvPath)

	if os.IsNotExist(err) {

		log.Fatal("The schema for reversion doesn't exist")
	}

	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := revertFilesFromCSV(f); err != nil {
		return err
	}

	return nil
}

func revertFilesFromCSV(f *os.File) error {
	csvReader := csv.NewReader(f)

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(rec) != 2 {
			log.Fatalf("unexpected number of columns in CSV: got %d, expected 2", len(rec))
		}
		// /home/aura/desktop/nameofthefile
		configPath := rec[0]
		//new or the backup folder path
		backupPath := rec[1]

		if backupPath == "new" {
			err = os.RemoveAll(configPath)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Removed:  ", configPath)

			continue
		} else {
			_, err = os.Stat(backupPath)
			if os.IsNotExist(err) {
				fmt.Printf("Source file does not exist: %s\n", backupPath)
				continue
			} else if err != nil {
				continue
			}

		}

		_, err = os.Stat(configPath)
		if !os.IsNotExist(err) {
			fmt.Println(configPath)
			err = os.RemoveAll(configPath)
			if err != nil {
				fmt.Printf("Error removing destination file %s: %v\n", configPath, err)
				continue
			}
		}

		err = os.Rename(backupPath, configPath)
		if err != nil {
			fmt.Printf("Error renaming file from %s to %s: %v\n", backupPath, configPath, err)
			continue
		}

		fmt.Printf("File reverted: %s to %s\n", backupPath, configPath)
	}

	return nil
}

func chooseBackupVersion(options []os.DirEntry) (os.DirEntry, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Choose the backup version\n")
		for i, dir := range options {
			prettyName, err := time.Parse("20060102150405", dir.Name())
			if err != nil {
				log.Fatal("You must have modified one of the directories. Don't do that", err)
			}
			cyan := "\033[0;36m"
			resetColor := "\033[0m"
			fmt.Printf("%d: %s%s%s\n", i+1, cyan, prettyName.Format("January 2, 2006 15:04:05"), resetColor)
		}

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		choice, err := strconv.Atoi(text)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println("Invalid choice, please choose a valid option.")
			continue
		}

		return options[choice-1], nil
	}
}
