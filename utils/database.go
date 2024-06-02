package utils

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
type FileProps struct{
		*gorm.Model
		DirPath string
		FilePath string
		Hash []byte
		WasChanged bool

	}


func InitDb() (*gorm.DB, error) {



db_path, err := Checkdir("trip_db.sqlite3")
	

	if err != nil {
		fmt.Println("Can't create a directroy ",err)
		return nil, err
	}

  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
		slog.Error("Can't open the db",err)
		return nil, err
  }
	return db, nil
}
