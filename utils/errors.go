package utils

import (
	"fmt"
	"os"
)
func HandleFileErr(action string, err error, file *os.File) error {
	if err != nil {
		fmt.Printf("Error %s: %v\n", action, err)
	}
	if file != nil {
		file.Close()        
		os.Remove(file.Name()) 
	}
	return err
}
