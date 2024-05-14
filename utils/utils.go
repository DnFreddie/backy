package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"encoding/json"


)

const (
	LOG_DIR = ".user_log"
)

func CheckJson() error {

	err := checkdir()
	if err != nil {

		return err
	}

	return nil
}










func Add_dir(paths *[]string) ([]string, error) {
    var newPaths []string 

    for _, p := range *paths {
        new_path, err := filepath.Abs(p)

        if err != nil {
            fmt.Println("Error getting absolute path for:", p)
            return nil, err
        }

        _, err = os.Stat(new_path)
        if os.IsNotExist(err) {
            fmt.Println("File does not exist:", p)
            continue 
        } else if err != nil {
            fmt.Println("Error checking file:", err)
            return nil, err
        }

        newPaths = append(newPaths, new_path)
    }

    return newPaths, nil
}



func checkdir() error {
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Can't get the user")
		return err
	}
	home_dir := user.HomeDir

	fmt.Printf(home_dir)
	log_dir := path.Join(home_dir, LOG_DIR)
	err = os.MkdirAll(log_dir, 0700)

	if err != nil {
		fmt.Printf("Cant create the %v", LOG_DIR)
		return err
	}

	fmt.Println("The log dir has been created", log_dir)
	return nil

}

func Jsonyfie(FDirs []string) error {
    b_Record := Brecord{
        Category: "test",
        FDirs:     FDirs,
        LMod:     "2024-05-14",
    }


	jr,err :=  json.Marshal(b_Record)
	if err != nil{
		fmt.Println("Can't marchal the recorde ",b_Record)
		return err

	}
	    fmt.Println(string(jr))
err=  os.WriteFile("test_file.json",jr,0666) 

	if err != nil{
		fmt.Println("can't write to a file",err)
		return err}



    return nil
}

