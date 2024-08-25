package dot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"github.com/DnFreddie/backy/utils"
)


func (r *Repo) Delete(path string) error {
	var choice string
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Are you sure you want to remove the backup? (y/n)")
		text, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(text)

		if choice != "y" && choice != "n" {
			fmt.Println("Invalid choice. Please enter 'y' or 'n'.")
			continue
		}

		break 
	}

	if choice == "y" {
		err := os.RemoveAll(path)
		if err != nil {
			slog.Error("Failed to remvoe backup %w","err",err)
			return err
		}
		fmt.Println("Backup removed successfully.")
	} else {
		fmt.Println("Backup removal canceled.")
	}

	return nil
}
// func (t *Repo) Revert() error {
// 	backupPath, err := utils.Checkdir(DOTS, false)
// 	if err != nil {
// 		return fmt.Errorf("failed to create backup: %w", err)
// 	}

// 	dirs, err := os.ReadDir(backupPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read directories: %w", err)
// 	}
// 	slog.Debug("Directories read:", "dirs", dirs)

// 	bV, err := utils.ChooseBackupVersion(dirs)

// 	if err != nil {
// 		return fmt.Errorf("failed to choose the version: %w", err)
// 	}

// 	chosenVersion := path.Join(DOTS, bV)
// 	r, err := processReversion(chosenVersion)
// 	if err != nil {
// 		return err
// 	}

	// for _, d := range *r.Dots {

	// 	if err := d.cleanLinks(); err != nil {
	// 		return err
	// 	}

	// 	err := os.RemoveAll(path.Join(backupPath, "dotfiles", bV))
	// 	if err != nil {
	// 		return &utils.UserError{FPath: path.Join(backupPath, "dotfiles", bV), Err: err}
	// 	}

	// }
	// return nil
// }

func (r *Repo) ExecuteReversion()error{
	slog.Debug("Started Revision")
	if len(*r.Dots) ==0 {
		
		return  fmt.Errorf("No files to revert")
	}
	for _, d := range *r.Dots {
		if err := d.cleanLinks(); err != nil {
			slog.Error("Failed to clean symlink","err",err)
			d.Failed =  err.Error()
		}
	}

return nil
}

// Reades the back_conf and Revert the schema to the previous state
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

	//If old (meaning it existed before)
	src := filepath.Join(d.Repo.BackupLocation, d.Location)
	if err := utils.Copy(src, d.Symlink); err != nil {
		slog.Error("Failed to copy from source to symlink:", "src", src, "symlink", d.Symlink, "error", err)
		return fmt.Errorf("failed to copy from %s to %s: %w", src, d.Symlink, err)
	}

	slog.Info("Successfully cleaned and updated symlink:", "symlink", d.Symlink, "source", src)
	return nil
}

func(r *Repo) ProcessRevision(chosenPath string) error {
	csvPath := path.Join(chosenPath, utils.SCHEMA_JSON)

	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		log.Fatal(&utils.UserError{Err: err, FPath: csvPath})
	}

	f, err := os.Open(csvPath)
	if err != nil {
		slog.Error("Failed to open file", "path", csvPath)
		return  err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(r); err != nil {
		log.Fatal(&utils.UserError{FPath: csvPath, Err: err})
		return  err
	}

	///Unmarchaling get's rid of the pointer to repo this
	//reassign it
	if r.Dots != nil {
		for i := range *r.Dots {
			(*r.Dots)[i].Repo = r
		}
	}
	return  nil
}
