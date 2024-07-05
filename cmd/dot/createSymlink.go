package dot

import (
	"encoding/csv"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"os"
	"path"
	"time"
)
func createTempBack(source string, backupDir string, csvF *csv.Writer, sourceAbs string, newDest string) (bool, error) {

	_, err := os.Stat(source)

	if os.IsNotExist(err) {
		return false, nil
	}

	dest := path.Join(backupDir, path.Base(source))
	fmt.Println("Already exist", path.Base(dest))
	err = os.Rename(source, dest)

	if err != nil {
		return false, err
	}
	data := [][]string{
		{source, dest},
	}
	err = csvF.WriteAll(data)

	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	err = os.Symlink(sourceAbs, newDest)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, nil

}

func createSymlink(dotfiles []Dotfile, source string) error {

	targetPath, err := utils.GetUser(TARGET)
	if err != nil {

		return err
	}

	nowT := time.Now().Format("20060102150405")

	backupDir, err := utils.Checkdir(path.Join(BACK_CONF, nowT), false)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(path.Join(backupDir, utils.SCHEMA_CSV))
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(f)

	for _, f := range dotfiles {

		if f.IsEx {
			symlinkPath := f.Location.Name()
			sourceAbs := path.Join(source, symlinkPath)
			dest := path.Join(targetPath, symlinkPath)

			wasCreated, err := createTempBack(dest, backupDir, writer, sourceAbs, dest)
			if err != nil {
				fmt.Println("Error creating temporary backup:", err)
				return err
			}

			if !wasCreated {

				err = os.Symlink(sourceAbs, dest)

				if err != nil {
					fmt.Println("Failed to create symlink:", err)
					return err

				}
				fmt.Println("Created symilnk", path.Base(dest))

				data := [][]string{
					{dest, "new"},
				}

				err = writer.WriteAll(data)

				if err != nil {
					return err
				}

			}
		}
	}

	return nil

}


