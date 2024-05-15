package utils

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
)


const (
	LOG_DIR   = ".user_log"
	JSON_PATH = "test_file.json"
)


func Checkdir() {
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Can't get the user")
		fmt.Println(err)
		os.Exit(1)
	}
	home_dir := user.HomeDir

	log_dir := path.Join(home_dir, LOG_DIR)
	err = os.MkdirAll(log_dir, 0700)

	if err != nil {
		fmt.Printf("Cant create the %v", LOG_DIR)
		fmt.Println(err)
		os.Exit(1)
	}

}

func ScanDir(dir_path string) ([]fs.DirEntry, error) {
	conf_dir, err := getUser(dir_path)

	if err != nil {
		fmt.Println("Scan failed: ", err)
		return nil, err
	}

	files, err := os.ReadDir(conf_dir)
	if err != nil {

		fmt.Printf("Doesnt have permissons to read {%v}", err)
	}

	return files, nil
}

// /Needd the path
func getUser(p string) (string, error) {
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

