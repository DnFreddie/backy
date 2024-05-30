package utils

import (
	"log/slog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {

db_path, err := Checkdir("trip_db.sqlite3")

	if err != nil {
		return nil, err
	}


  db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
  if err != nil {
		slog.Error("Can't open the db",err)
		return nil, err
  }
	return db, nil
}
