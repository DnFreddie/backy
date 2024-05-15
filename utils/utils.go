package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	LOG_DIR = ".user_log"
	JSON_PATH = "test_file.json"
)

func CheckJson() error {

	err := checkdir()
	if err != nil {

		return err
	}

	return nil
}

func Add_dir(paths *[]string) ([]string, error) {
	var newPaths []string

	for _, p := range *paths {
		new_path, err := filepath.Abs(p)

		if err != nil {
			fmt.Println("Error getting absolute path for:", p)
			return nil, err
		}

		_, err = os.Stat(new_path)
		if os.IsNotExist(err) {
			fmt.Println("File does not exist:", p)
			continue
		} else if err != nil {
			fmt.Println("Error checking file:", err)
			return nil, err
		}

		newPaths = append(newPaths, new_path)
	}

	return newPaths, nil
}

func checkdir() error {
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Can't get the user")
		return err
	}
	home_dir := user.HomeDir

	fmt.Printf(home_dir)
	log_dir := path.Join(home_dir, LOG_DIR)
	err = os.MkdirAll(log_dir, 0700)

	if err != nil {
		fmt.Printf("Cant create the %v", LOG_DIR)
		return err
	}

	fmt.Println("The log dir has been created", log_dir)
	return nil

}

func Jsonyfie(FDirs []string) error {
	b_Record := Brecord{
		Category: "test",
		FDirs:    FDirs,
		LMod:     "2024-05-14",
	}

	_, err := os.Stat("test_file.json")
	if os.IsNotExist(err) {
		fmt.Println("File does not exist creating one ...")
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

	jsonData, err := ReadJson(JSON_PATH)
	if err != nil {
		return err
	}

	// Append new record to existing data
	jsonData = append(jsonData, b_Record)

	jr, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Can't marshal the records ", jsonData)
		return err
	}

	err = os.WriteFile("test_file.json", jr, 0666)
	if err != nil {
		fmt.Println("Can't write to a file", err)
		return err
	}

	return nil
}

func ReadJson(jsonPath string) ([]Brecord, error) {
	var bR []Brecord
	f, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println("Can't read the file")
		return nil, err
	}

	fmt.Println(string(f))
	err = json.Unmarshal(f, &bR)
	if err != nil {
		fmt.Println("Can't unmarshal the records ", err)
		return nil, err
	}

	return bR, nil
}


