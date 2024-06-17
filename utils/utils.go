package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	LOG_DIR   = ".user_log"
)

func Checkdir(fPath string, file bool) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("can't get the user: %v", err)
	}
	homeDir := user.HomeDir

	logDir := path.Join(homeDir, LOG_DIR)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("can't create the directory %v: %v", LOG_DIR, err)
	}

	requestedF := path.Join(logDir, fPath)
	_, err = os.Stat(requestedF)
	if os.IsNotExist(err) {

		if file {
			f, err := os.Create(requestedF)
			if err != nil {
				return "", fmt.Errorf("can't create the file %v due to: %v", requestedF, err)
			}
			defer f.Close()
			return requestedF, nil
		} else {
			err = os.MkdirAll(requestedF, os.ModePerm)
			if err != nil {
				return "", fmt.Errorf("can't create the file %v due to: %v", requestedF, err)
			}

		}

	} else if err != nil {
		return "", fmt.Errorf("error checking for file %v: %v", requestedF, err)
	}

	return requestedF, nil
}

func ScanDir(dir_path string) ([]fs.DirEntry, error) {
	_, err := os.Stat(dir_path)

	if os.IsNotExist(err) || err != nil {
		fmt.Println("Scan failed: ", err)
	}

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Printf("Doesnt have permissons to read {%v}", err)
		return nil, err
	}

	return files, nil
}

// /Returns the joined path of the target and the user dir
func GetUser(p string) (string, error) {
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Can't get the user")
		fmt.Println(err)
		os.Exit(1)
	}
	home_dir := user.HomeDir
	joined_path := path.Join(home_dir, p)

	return joined_path, nil

}



func ReadJson[T any](jsonPath string, records *[]T) error {
	f, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println("Can't read the file:", err)
		return err
	}
	//TODO makei so this returs error and the ccheck for it's type

	if len(f) == 0 {
		return nil
	}

	err = json.Unmarshal(f, records)
	if err != nil {
		fmt.Println("Can't unmarshal the records:", err)
		return err
	}

	return nil
}
func MakeAbsoulute(fPath string) (string, error) {
	var dest string
	if !filepath.IsAbs(fPath) {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Therse smth wrong with this directroy check perrmisons ")
		}
		dest = filepath.Join(pwd, fPath)

	} else {
		dest = fPath
	}
	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		return "", err
	}

	return dest, nil
}
