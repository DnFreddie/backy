package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/DnFreddie/backy/utils"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Reades the back_conf and Revert the schema to the previous state
func (t *Repo)Revert() error {

	backupPath, err := utils.Checkdir(BACK_CONF, false)
	if err != nil {
		log.Fatal("Failed to create backup")
	}

	//TODO! Fix this read config 
	dirs, err := os.ReadDir(path.Join(backupPath,"dotfiles"))
	slog.Debug("This is ", "dirs",dirs)

	bV, err := chooseBackupVersion(dirs)
	if err != nil {
		log.Fatal("Failed to chose the veriosn ", err)
	}
	r, err := processReversion(path.Join(backupPath,"dotfiles", bV))

	if err != nil {

		log.Fatal(err)

	}
	for _, dot := range *r.Dots {
		


		dot.cleanLinks()
	}
	return nil

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

func processReversion(chosenPath string) (Repo, error) {
	r := &Repo{}
	csvPath := path.Join(chosenPath, utils.SCHEMA_JSON)

	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		log.Fatal(&utils.UserError{Err: err, FPath: csvPath})
	}

	f, err := os.Open(csvPath)
	if err != nil {
		slog.Error("Failed to open file", "path", csvPath)
		return Repo{}, err
	}
	defer f.Close() 

	if err := json.NewDecoder(f).Decode(r); err != nil {
		log.Fatal(&utils.UserError{FPath: csvPath, Err: err})
		return Repo{}, err
	}

	return *r, nil
}



