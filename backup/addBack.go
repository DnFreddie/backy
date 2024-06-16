package backup

import (
	"encoding/json"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"os"
)

func Add_dir(paths *[]string) ([]string, error) {
	var newPaths []string

	for _, p := range *paths {
		new_path, err := utils.MakeAbsoulute(p)
		if err != nil {
			fmt.Println(p,"Doesn't exist")

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

func Jsonyfie(FDirs []string) error {
	b_Record := Brecord{
		Category: "test",
		FDirs:    FDirs,
		LMod:     "2024-05-14",
	}

	_, err := os.Stat("test_file.json")
	if os.IsNotExist(err) {
		fmt.Printf("/n File test_file.json  does not exist creating one ...")
		b_RecordArray := []Brecord{b_Record}
		jr, err := json.Marshal(b_RecordArray)
		if err != nil {
			fmt.Println("Can't marshal the record ", b_Record)
			return err
		}

		err = os.WriteFile("test_file.json", jr, 0666)
		if err != nil {
			fmt.Println("Can't write to a file", err)
			return err
		}
	}

	var Record []Brecord
	err = utils.ReadJson(utils.JSON_PATH, &Record)
	if err != nil {
		return err
	}

	Record = append(Record, b_Record)

	jr, err := json.Marshal(Record)
	if err != nil {
		fmt.Println("Can't marshal the records ", Record)
		return err
	}

	err = os.WriteFile("test_file.json", jr, 0666)
	if err != nil {
		fmt.Println("Can't write to a file", err)
		return err
	}

	return nil
}
