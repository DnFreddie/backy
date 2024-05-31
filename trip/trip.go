package trip

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sync"

	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

type FileChanged struct {
	AbPath     string
	Hash       []byte
	WasChanged bool
}

func HashHash(b []byte, fPath string) FileChanged {
	hash := sha256.Sum256(b)

	testChanged := FileChanged{
		AbPath:     fPath,
		Hash:       hash[:],

		WasChanged: false,
	}

	return testChanged
}

func readFilesAsync(filePath string) (FileChanged,error){

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist:", filePath)
		return FileChanged{},nil
	}

	content, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Can't read the file permissions issues", err)
		return FileChanged{}, err
	}

	fChanged := HashHash(content, filePath)


	return fChanged ,nil
}

func TripCommand(fPath string, db *gorm.DB, ch chan utils.FileProps) error {
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
                err := TripCommand(jPath, db, ch)
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

func TripMain(fPath string) error {
    ch := make(chan utils.FileProps)
    done := make(chan bool)

    db, err := utils.InitDb()
    if err != nil {
        return err
    }

    err = db.AutoMigrate(&utils.FileProps{})
    if err != nil {
        return err
    }

    go func() {
        TripCommand(fPath, db, ch)
        close(ch)  
    }()

    go processxatches(ch, done, db)

    <-done


    var count int64
    if err := db.Model(&utils.FileProps{}).Count(&count).Error; err != nil {
        return err
    }
    fmt.Println("Total count of FileProps:", count)

    return nil
}






func processxatches(ch chan utils.FileProps, done chan<- bool, db *gorm.DB) {
    for {

        select {
        case item, ok := <-ch:
            if !ok {
		fmt.Println("im done ")
                done <- true
                return
            }
			var count int64
			db.Model(&utils.FileProps{}).Where("file_path = ?", item.FilePath).Count(&count)
			 if count == 0  {
                    fmt.Println(item.FilePath, "does not exist in the file_props table.")
					fmt.Println(item.DirPath)
					count=0
					
			}

        }
    }
}







func processBatches(ch chan utils.FileProps, done chan<- bool, db *gorm.DB) {
    var batch []utils.FileProps
    var count int
    for {
        select {
        case item, ok := <-ch:
            if !ok {
                if len(batch) > 0 {
                    db.CreateInBatches(batch, len(batch))
                    fmt.Println("Processed count during execution:", count)
                }
                done <- true
                return
            }
            count++
            batch = append(batch, item)
            if len(batch) == 200 {
                db.CreateInBatches(batch, len(batch))

                fmt.Println("Processed count during execution:", count)
				batch = batch[:0] // clear the batch

            }
        }
    }
}

