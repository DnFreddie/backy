package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func revert() {

	dirs, err := os.ReadDir("/home/rocky/.user_log/back_conf/")
	bV, err := chooseBackupVersion(dirs)
	if err != nil {
		log.Fatal("Failed to chose the veriosn ", err)
	}
	r, err := processReversion(path.Join("/home/rocky/.user_log/back_conf/", bV))

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(r)

	for _, i := range *r.Dots {
		fmt.Printf("Moving %v to %v is new %v\n", i.Symlink, i.Absolute, i.New)
	}

}

func (d *Dotfile) cleanLinks() error {
	if d.Symlink == "" {
		return nil
	}

	assertSymlink, err := os.Lstat(d.Symlink)
	if err != nil {
		return fmt.Errorf("failed to stat symlink %s: %w", d.Symlink, err)
	}

	if assertSymlink.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("this is not a symlink: %s, you might have changed it, skipping", d.Symlink)
	}

	if d.New {
		fmt.Println("Removing the symlink:", d.Symlink)
		if err := os.Remove(d.Symlink); err != nil {
			return fmt.Errorf("failed to remove symlink %s: %w", d.Symlink, err)
		}
		return nil
	}

	if err := os.Remove(d.Symlink); err != nil {
		return fmt.Errorf("failed to remove symlink %s: %w", d.Symlink, err)
	}

	src := filepath.Join(d.Repo.BackupLocation, d.Location)
	if err := utils.Copy(src, d.Symlink); err != nil {
		return fmt.Errorf("failed to copy from %s to %s: %w", src, d.Symlink, err)
	}

	return nil
}

func chooseBackupVersion(options []os.DirEntry) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose the backup version")

		for i := len(options) - 1; i >= 0; i-- {
			dir := options[i]
			prettyName, err := time.Parse("20060102150405", dir.Name())
			if err != nil {
				log.Fatalf("You must have modified one of the directories. Don't do that: %v", err)
			}
			fmt.Printf("%d: %s%s%s\n", i+1, Cyan, prettyName.Format("January 2, 2006 15:04:05"), Reset)
		}

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		choice, err := strconv.Atoi(text)
		if err != nil || choice < 1 || choice > len(options) {
			fmt.Println("Invalid choice, please choose a valid option.")
			continue
		}

		return options[len(options)-choice].Name(), nil
	}
}

func processReversion(chosenPath string) (Repo, error) {
	r := &Repo{}
	csvPath := path.Join(chosenPath, "backy_schema.josn")
	fmt.Println(csvPath)

	_, err := os.Stat(csvPath)

	if os.IsNotExist(err) {

		log.Fatal("The schema for reversion doesn't exist")
	}

	f, err := os.Open(csvPath)
	if err != nil {

	}
	defer f.Close()

	bytes, err := io.ReadAll(f)

	json.Unmarshal(bytes, r)

	return *r, nil
}
