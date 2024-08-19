package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Reades the back_conf and revert the schema to the previous state
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
		slog.Info("No symlink to clean, skipping.")
		return nil
	}

	assertSymlink, err := os.Lstat(d.Symlink)
	if err != nil {
		slog.Error("Failed to stat symlink:", "symlink", d.Symlink, "error", err)
		return err
	}

	if assertSymlink.Mode()&os.ModeSymlink == 0 {
		return &utils.UserError{
			Err:   fmt.Errorf("not a symlink"),
			FPath: d.Symlink,
		}
	}

	if d.New {
		slog.Info("Removing the symlink:", "symlink", d.Symlink)
		if err := os.Remove(d.Symlink); err != nil {
			slog.Error("Failed to remove symlink:", "symlink", d.Symlink, "error", err)
			return fmt.Errorf("failed to remove symlink %s: %w", d.Symlink, err)
		}
		return nil
	}

	if err := os.Remove(d.Symlink); err != nil {
		slog.Error("Failed to remove symlink:", "symlink", d.Symlink, "error", err)
		return fmt.Errorf("failed to remove symlink %s: %w", d.Symlink, err)
	}

	src := filepath.Join(d.Repo.BackupLocation, d.Location)
	if err := utils.Copy(src, d.Symlink); err != nil {
		slog.Error("Failed to copy from source to symlink:", "src", src, "symlink", d.Symlink, "error", err)
		return fmt.Errorf("failed to copy from %s to %s: %w", src, d.Symlink, err)
	}

	slog.Info("Successfully cleaned and updated symlink:", "symlink", d.Symlink, "source", src)
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
				log.Fatal(&utils.UserError{Err: err, FPath: dir.Name()})
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

// Reades the Schema and provides the formmer Repo struckt based on that
func processReversion(chosenPath string) (Repo, error) {
	r := &Repo{}
	csvPath := path.Join(chosenPath, "backy_schema.josn")
	fmt.Println(csvPath)

	_, err := os.Stat(csvPath)

	if os.IsNotExist(err) {

		log.Fatal(&utils.UserError{Err: err, FPath: csvPath})
	}

	f, err := os.Open(csvPath)
	defer f.Close()
	if err != nil {
		slog.Error("Faield to open file", "path", csvPath)
		return Repo{}, err

	}

	bytes, err := io.ReadAll(f)

	err = json.Unmarshal(bytes, r)
	log.Fatal(utils.UserError{FPath: csvPath, Err: err})

	return *r, nil
}
