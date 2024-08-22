
package common
import "log"

type Backup interface {
	SaveSchema() error
	PrintRaport() 
	Revert() error
}

func SaveStructure (backup Backup) error {
	if err := backup.SaveSchema(); err != nil {
		log.Println("Error saving repo schema:", err)
		return err
	}
	backup.PrintRaport() 
	return nil
}
