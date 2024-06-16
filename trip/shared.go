package trip

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"path"
	"sync"
)

const DB_PATH = "trip_db.sqlite3"

type FileChanged struct {
	AbPath     string
	Hash       []byte
	WasChanged bool
}

func ScanRecursivly(fPath string, db *gorm.DB, ch chan utils.FileProps) error {
	fds, err := os.ReadDir(fPath)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, fd := range fds {
		jPath := path.Join(fPath, fd.Name())
		if fd.IsDir() {
			wg.Add(1)
			go func(jPath string, db *gorm.DB, ch chan utils.FileProps) {
				defer wg.Done()
				err := ScanRecursivly(jPath, db, ch)
				if err != nil {
					fmt.Printf("Error in directory %s: %v\n", jPath, err)
					return
				}
			}(jPath, db, ch)
		} else {

			has, err := readFilesAsync(jPath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", jPath, err)
				continue
			}
			newProps := utils.FileProps{
				DirPath:    fPath,
				FilePath:   has.AbPath,
				Hash:       has.Hash,
				WasChanged: has.WasChanged,
			}
			ch <- newProps
		}
	}
	wg.Wait()
	return nil
}

func readFilesAsync(filePath string) (FileChanged, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Can't read the File propably  does not exist:", filePath)
		return FileChanged{}, nil
	}

	content, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Can't read the file permissions issues", err)
		return FileChanged{}, err
	}

	fChanged := hashHash(content, filePath)

	return fChanged, nil
}

func hashHash(b []byte, fPath string) FileChanged {
	hash := sha256.Sum256(b)
	testChanged := FileChanged{
		AbPath:     fPath,
		Hash:       hash[:],
		WasChanged: false,
	}

	return testChanged
}

type ConfigPath struct {
	Fpath string `json:"fpath"`
}

func CreateConfig(scanPath string) (bool, error) {

	confP, err := utils.Checkdir("scan_paths.json", true)
	var existed []ConfigPath

	err = utils.ReadJson(confP, &existed)
	if err != nil {
		fmt.Println("Can't read json")
		return true, err
	}

	var exist bool
	for _, p := range existed {

		if p.Fpath == scanPath {
			exist = true
			break
		}
	}

	if !exist {
		existed = append(existed, ConfigPath{Fpath: scanPath})
		bytes, err := json.Marshal(&existed)
		if err != nil {
			return false, err
		}
		os.WriteFile(confP, bytes, 0666)
	}

	return exist, nil
}
