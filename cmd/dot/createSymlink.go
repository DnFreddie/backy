package dot

import (
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

func (r *Repo) createRaport(){
	if r.BackupLocation == ""{
		log.Fatal("Can't find the backup location")

	}
	var db *gorm.DB
	db,err:=  utils.InitDb("repos.sql",Repo{})
	db,err=  utils.InitDb("repos.sql",Dotfile{})
	if err != nil {
		log.Fatal(err)
	}


	db.Create(&r)
	for _,i := range *r.Dots{
		db.Create(&i)

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
			d.Failed =&strError
				return

			}

		}
		err := os.Symlink(d.Symlink, d.Repo.Absolute)

		if err != nil {
			strError := err.Error()
			d.Failed =&strError
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
