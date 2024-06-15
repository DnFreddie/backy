package revert

import (
	"fmt"
	"os"
	"path"
	"github.com/DnFreddie/backy/dot"
	"github.com/DnFreddie/backy/utils"
)

func deletBackup() error {

	confDir, err := utils.Checkdir(dot.BACK_CONF, false)
	if err != nil {
		return err
	}
	options, err := os.ReadDir(confDir)
	if err != nil {
		return err
	}

	if len(options) < 1 {
		fmt.Println("No backups to chose")
		return nil
	}

	revDir, err := chooseBackupVersion(options)
	if err != nil {
		return err
	}
	chosenPath := path.Join(confDir, revDir.Name())
	err = os.RemoveAll(chosenPath)

	if err != nil {
		return err
	}
	fmt.Printf("\n %v has been removed\n", revDir.Name())
	return nil

}
