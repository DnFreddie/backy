package common

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/DnFreddie/backy/utils"
)

// ProcessRevision defines a function type for processing a revision.

// ChooseBackupVersion defines a function type for choosing a backup version.
type ChooseBackupVersion func([]os.DirEntry) (string, error)

// BackupManager defines methods for managing backups.
type BackupManager interface {
	SaveSchema() error
	PrintReport() 
	ProcessRevision(string)error
	ExecuteReversion() error
}

// Revert reverts to a chosen backup version.
func Revert(b BackupManager,backupPath string) error {

	dirs, err := os.ReadDir(backupPath)
	if err != nil {
		return fmt.Errorf("error reading backup directory: %w", err)
	}

	bV, err := utils.ChooseBackupVersion(dirs)
	if err != nil {
		return fmt.Errorf("error choosing backup version: %w", err)
	}

	choice := (path.Join(backupPath,bV))
	if err := b.ProcessRevision(choice); err != nil {
		return fmt.Errorf("error processing revision: %w", err)
	}

	err = b.ExecuteReversion()

	if  err!= nil{
		slog.Error("Failed to ExecuteReversion","err",err)
		return err
	}

		if err := os.RemoveAll(path.Join(backupPath,bV )); err != nil {
		return fmt.Errorf("error removing backup directory: %w", err)
	}

	return nil
}

// SaveStructure saves the backup schema and prints a report.
func SaveStructure(b BackupManager) error {
	if err := b.SaveSchema(); err != nil {
		log.Println("Error saving repo schema:", err)
		return err
	}
	 b.PrintReport()	
	return nil
}

func DeleteBackup(b BackupManager, bPath string) error {
	dir, err := os.ReadDir(bPath)
	if err != nil {
		return fmt.Errorf("error reading backup directory: %w", err)
	}

	backupVersion, err := utils.ChooseBackupVersion(dir) 
	if err != nil {
		log.Println("Error choosing backup version:", err)
		return fmt.Errorf("error choosing backup version: %w", err)
	}

	if err := os.RemoveAll(path.Join(bPath, backupVersion)); err != nil {
		return fmt.Errorf("error removing backup directory: %w", err)
	}
	return nil
}

