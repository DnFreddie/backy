package trip

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"

	"github.com/DnFreddie/backy/utils"
	"log/slog"

	"os"

	"gorm.io/gorm"
)

type Status string

const (
	NotExist Status = "Not Exists"
	Changed  Status = "Changed"
)

type Compared struct {
	Status
	FilePath string
	Dir      string
}

func checkSumFiles(ch chan utils.FileProps, db *gorm.DB, fchan chan Compared) {
	for item := range ch {
		user := utils.FileProps{}
		comperdItem := Compared{
			FilePath: item.FilePath,
			Dir:      item.DirPath,
		}

		db.Where("file_path = ?", item.FilePath).Find(&user)

		if reflect.ValueOf(user).IsZero() {
			comperdItem.Status = NotExist
			fchan <- comperdItem
		} else {
			compersion := compere(user, item)
			if compersion != "" {
				comperdItem.Status = Changed
				fchan <- comperdItem
			}

		}
	}

}

func compere(saved utils.FileProps, newI utils.FileProps) string {

	if !bytes.Equal(saved.Hash, newI.Hash) {
		fmt.Println("the hasesh doens't work ", saved.FilePath)
		return saved.FilePath
	}

	return ""

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
