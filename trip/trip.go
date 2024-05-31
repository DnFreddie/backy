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
		AbPath: fPath,
		Hash:   hash[:],
		WasChanged: false,
	}

	return testChanged
}

func readFilesAsync(filePath string) (FileChanged, error) {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist:", filePath)
		return FileChanged{}, nil
	}

	content, err := os.ReadFile(filePath)

	if err != nil {
		slog.Error("Can't read the file permissions issues", err)
		return FileChanged{}, err
	}

	fChanged := HashHash(content, filePath)

	return fChanged, nil
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
    var wg sync.WaitGroup

    db, err := utils.InitDb()
    if err != nil {
        return err
    }

    numWorkers := 200
    wg.Add(numWorkers)

    if !db.Migrator().HasTable(&utils.FileProps{}) {
        err = db.AutoMigrate(&utils.FileProps{})
	
	if err != nil {
		return err
	}

    go func() {
        TripCommand(fPath, db, ch)
        close(ch)
    }()

     for i := 0; i < numWorkers; i++ {
            go processBatches(ch, &wg, db)
    }

       
    }else{
    go func() {
        TripCommand(fPath, db, ch)
        close(ch)
    }()

     for i := 0; i < numWorkers; i++ {
            go processBatches(ch, &wg, db)
    }




	}

    wg.Wait()
    var count int64
    if err := db.Model(&utils.FileProps{}).Count(&count).Error; err != nil {
        return err
    }
    fmt.Println("Total count of FileProps:", count)

    return nil
}


func processxatches(ch chan utils.FileProps, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	for item := range ch {
		var count int64
		db.Model(&utils.FileProps{}).Where("file_path = ?", item.FilePath).Count(&count)
		if count == 0 {
			fmt.Println(item.FilePath, "does not exist in the file_props table.")
			fmt.Println(item.DirPath)
		}
	}
}

func processBatches(ch chan utils.FileProps, wg *sync.WaitGroup, db *gorm.DB) {
    defer wg.Done()
    var batch []utils.FileProps
    var count int

    for item := range ch {
        count++
        batch = append(batch, item)
        if len(batch) == 300 {
            db.CreateInBatches(batch, len(batch))
            fmt.Println("Processed count during execution:", count)
            batch = batch[:0] 
        }
    }

    if len(batch) > 0 {
        db.CreateInBatches(batch, len(batch))
        fmt.Println("Processed count during execution:", count)
    }
}
