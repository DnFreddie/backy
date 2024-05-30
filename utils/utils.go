package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path"
)

const (
	LOG_DIR   = ".user_log"
	JSON_PATH = "test_file.json"
)

func Checkdir(fPath string) (string,error) {
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

	requestedF := path.Join(log_dir, fPath)
	_, err = os.Stat(requestedF)

	if os.IsNotExist(err) {
		f, err := os.Create(fPath)
		if err != nil {
			fmt.Printf("Can't create the file %v due to %v\n", fPath, err)
			return "",err
		}
		defer f.Close()
	} else if err != nil {
		return "",nil
	}

	return "",nil
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

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func CopyDir(src string, dst string) error {
	var err error
	var fds []fs.DirEntry
	var srcinfo fs.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func ReadJson[T any](jsonPath string, unmarshalS *T) ([]T, error) {
	var records []T

	f, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println("Can't read the file:", err)
		return nil, err
	}

	err = json.Unmarshal(f, unmarshalS)
	if err != nil {
		fmt.Println("Can't unmarshal the records:", err)
		return nil, err
	}

	records = append(records, *unmarshalS)

	return records, nil
}
