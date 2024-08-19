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




type UserError struct {
	Err error
	FPath string 
}

func (u *UserError) Error() string {
	return fmt.Sprintf("user error: you changed something in this file; please check :%v. original error: %v", u.Err,u.FPath)
}

