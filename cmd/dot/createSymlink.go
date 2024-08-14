package dot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

func (r *Repo) createBackup() {
	nowT := time.Now().Format("20060102150405")
	backupDir := path.Join(BACK_CONF, nowT)
	backupPath, err := utils.Checkdir(backupDir, false)
	if err != nil {
		log.Fatal("Failed to create backup")
	}

	r.BackupLocation = backupPath

}

func (r *Repo) createDbRaport() {
	if r.BackupLocation == "" {
		log.Fatal("Can't find the backup location")
	}

	var db *gorm.DB
	db, dbInitErr := utils.InitDb("repos.sql", Repo{})
	if dbInitErr != nil {
		log.Println("Failed to create the database:", dbInitErr)
	}

	db, dbInitErr = utils.InitDb("repos.sql", Dotfile{})
	if dbInitErr != nil {
		log.Println("Failed to create the database:", dbInitErr)
	}

	if err := db.Create(&r).Error; err != nil {
		log.Println("Error creating repo in database:", err)
	}

	batchSize := 30
	if err := db.CreateInBatches(*r.Dots, batchSize).Error; err != nil {
		log.Println("Error creating dots in batches:", err)
	}

	jsonData, jsonErr := json.MarshalIndent(*r, "", "  ")
	if jsonErr != nil {
		log.Println("Failed to marshal data:", jsonErr)
	}

	schemaDest := path.Join(r.BackupLocation, "backy_schema.josn")
	file, writeErr := os.Create(schemaDest)
	if writeErr != nil {
		log.Println("Failed to create output.json:", writeErr)
	}

	if file != nil {
		defer file.Close()
		if _, writeErr := file.Write(jsonData); writeErr != nil {
			log.Println("Failed to write JSON data to file:", writeErr)
		}
	}

	if dbInitErr != nil && jsonErr != nil {
		panic("Both database initialization and JSON operations failed")
	}

	if writeErr != nil && dbInitErr != nil {
		log.Println("JSON Output:", string(jsonData))
	}
}

func (r *Repo) createCsrRaport() {
	if r.BackupLocation == "" {
		log.Fatal("Can't find the backup location")

	}
}

func (r *Repo) Link() {
	if r.Dots == nil {
		fmt.Println("No Files to link ")
		return

	}

	target, err := utils.GetUser(TARGET)

	if err != nil {
		log.Fatal("Failed to read the config", err)
	}
	r.createBackup()
	for _, dot := range *r.Dots {
		dot.IsExe()
		dot.createSymlink(target)
	}

}

func (d *Dotfile) createSymlink(target string) {
	if d.Executable {
		d.isNew(target)
		if !d.New {

			err := os.Rename(d.Symlink, d.Repo.BackupLocation)
			if err != nil {
				strError := err.Error()
				d.Failed = &strError
				return

			}

		}
		err := os.Symlink(d.Symlink, d.Repo.Absolute)

		if err != nil {
			strError := err.Error()
			d.Failed = &strError
			return
		}

	}

}

func (d *Dotfile) isNew(target string) {
	d.Symlink = path.Join(target, d.Location)
	_, err := os.Stat(d.Symlink)

	if os.IsNotExist(err) {
		d.New = true
	}

}
