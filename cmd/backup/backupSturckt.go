package backup

import "io/fs"

type Backup struct {

}
type Brecord struct {
	Parent *fs.DirEntry
	BackupPath string 
	Origin string
	DateAdded string
	Updaated string
	Hash string
}
func add (b Backup)(){


	
}

func (b *Backup) SaveSchema() error {
	return nil
}

func (b *Backup) PrintReport() {
}

func (b *Backup) ProcessRevision(revision string) error {
	return nil
}

func (b *Backup) ExecuteReversion() error {
	return nil
}

