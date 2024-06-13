package dot

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/DnFreddie/backy/utils"
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

	reader := bufio.NewReader(os.Stdin)
	var revDir fs.DirEntry

	for {
		fmt.Print("Choose the backup version\n")
		for i, dir := range options {

					
			prettyName, err := time.Parse("20060102150405", dir.Name())
			if err !=nil{
				log.Fatal("U must have modyfied one of the directories Dont'Do that",err)
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

		revDir = options[choice-1]
		break
	}

	chosenPath := path.Join(confDir, revDir.Name())

	csvPath := path.Join(chosenPath, "test.csv")
	_, err = os.Stat(csvPath)
	if os.IsNotExist(err) {
		return errors.New("The schema for revertion doesn't exist")
	}
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", rec)
	}

	return nil
}
