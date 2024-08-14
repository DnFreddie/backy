package utils

import (
	"fmt"
	"os"
)
func HandleFileErr(action string, err error, file *os.File) error {
	if err != nil {
		return fmt.Errorf("error %s: %v", action, err)
	}

	if file != nil {
		if removeErr := os.Remove(file.Name()); removeErr != nil {
			return fmt.Errorf("failed to remove file %s: %v. Please remove it manually.", file.Name(), removeErr)
		}
	}

	return nil
}

