package backup

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/DnFreddie/backy/utils"
)

func Add_dir(paths *[]string) ([]string, error) {
	var newPaths []string

	for _, p := range *paths {
		new_path, err := utils.MakeAbsoulute(p)
		if err != nil {
			fmt.Println(p, "Doesn't exist")

			continue
		}

		newPaths = append(newPaths, new_path)
	}

	return newPaths, nil
}

type Brecord struct {
	Category string   `json:"category"`
	FDirs    []string `json:"f_dir"`
	LMod     string   `json:"l_mod"`
	Changes  Mchanges `json:"changes"`
}

type Mchanges struct {
	Changed []string `json:"changed"`
	MTime   []string `json:"m_time"`
}

const (
	BACK_PATH = "backy_back.json"
)

func Jsonyfie(FDirs []string) error {
	b_Record := Brecord{
		Category: "test",
		FDirs:    FDirs,
		LMod:     "2024-05-14",
	}

	bp, err := utils.Checkdir(BACK_PATH, true)
	if err != nil {
		log.Fatal(err)
	}

	var Record []Brecord

	err = utils.ReadJson(bp, &Record)

	if err != nil {
		return err
	}

	Record = append(Record, b_Record)

	jr, err := json.Marshal(Record)
	if err != nil {
		fmt.Println("Can't marshal the records ", Record)
		return err
	}

	err = os.WriteFile(bp, jr, os.ModePerm)
	if err != nil {
		fmt.Println("Can't write to a file", err)
		return err
	}

	return nil
}
