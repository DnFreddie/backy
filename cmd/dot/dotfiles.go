package dot

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DnFreddie/backy/utils"
)

type Dotfile struct {
	ID         uint    `gorm:"primaryKey;autoIncrement" json:"-"`
	Location   string  `json:"location"`
	Executable bool    `json:"executable"`
	Symlink    string  `json:"symlink"`
	Absolute   string  `json:"ab_path"`
	Repo       *Repo   `gorm:"-" json:"-"`
	RepoID     *string `json:"-"`
	Ignored    bool    `json:"ignored"`
	New        bool    `json:"new"`
	Failed     *string `json:"failed,omitempty"`
}

func (d *Dotfile) createSymlink(target string)error {
	d.isNew(target)
	if !d.New {
		backup_dir := path.Join(d.Repo.BackupLocation, path.Base(d.Symlink))
		err := utils.Copy(d.Symlink, backup_dir)
		if err != nil {
			strError := err.Error()
			fmt.Println("Failed to rename: ", strError)
			d.Failed = &strError
			return err
		}
		err = os.RemoveAll(d.Symlink)
		if err != nil{
			log.Fatal("Failed to remove the symlink",d.Symlink,err)
		}
	}
	if d.Symlink == "" {
		err := fmt.Sprintf("the symlink path shouldn't be empty")
		log.Fatal(err)
		fmt.Println(err)
		d.Failed = &err
		return nil
	}
	err := os.Symlink(d.Absolute, d.Symlink)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}


func (d *Dotfile) isNew(target string) {
	d.Symlink = path.Join(target, d.Location)
	if d.Symlink == ""{
log.Fatal("Wtif is this ")
	}
	_, err := os.Stat(d.Symlink)
	if os.IsNotExist(err) {
		d.New = true
	} else {
		d.New = false
	}
}
func (d *Dotfile) IsExe() {
	re := regexp.MustCompile(`\.?(conf|rc)$`)
	cleanPath := re.ReplaceAllString(d.Location, "")
	d.Executable = isCmd(strings.TrimPrefix(cleanPath, "."))
}

func isCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	
	if err != nil {
		return false 
	}
	
	return true 
}

func (d *Dotfile) ignore(toIgnore *[]string) {
	d.Ignored = false

	for _, pattern := range *toIgnore {
		if match, _ := filepath.Match(pattern, d.Location); match {
			d.Ignored = true
			break
		}
	}
}
