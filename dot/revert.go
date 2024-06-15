package dot

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

func RevertBackups() error {
	confDir, err := utils.Checkdir(BACK_CONF, false)
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

	if err := processReversion(confDir, revDir); err != nil {
		return err
	}

	err = os.RemoveAll(confDir)
	if err != nil {
		log.Fatal("Have you changed permission? This shouldn't have happened")
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
			return fmt.Errorf("unexpected number of columns in CSV: got %d, expected 2", len(rec))
		}

		source := rec[1]
		dest := rec[0]

		_, err = os.Stat(source)
		if os.IsNotExist(err) {
			fmt.Printf("Source file does not exist: %s\n", source)
			continue
		} else if err != nil {
			return err
		}

		_, err = os.Stat(dest)
		if !os.IsNotExist(err) {
			err = os.Remove(dest)
			if err != nil {
				fmt.Printf("Error removing destination file %s: %v\n", dest, err)
				continue
			}
		}

		err = os.Rename(source, dest)
		if err != nil {
			fmt.Printf("Error renaming file from %s to %s: %v\n", source, dest, err)
			continue
		}

		fmt.Printf("File reverted: %s to %s\n", source, dest)
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

func processReversion(confDir string, revDir os.DirEntry) error {

	chosenPath := path.Join(confDir, revDir.Name())
	fmt.Println(revDir.Name())

	csvPath := path.Join(chosenPath, REVERT_CSV)
	_, err := os.Stat(csvPath)
	fmt.Println("this is the csv path ", csvPath)
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
