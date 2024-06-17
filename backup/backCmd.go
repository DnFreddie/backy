package backup

import "fmt"

func Add_command(args *[]string) error {
	paths, err := addDir(args)

	if err != nil {
		return err
	}
	if len(paths)==0{
		fmt.Println("Skipping... Nothing to add ")
		return nil
	}
	err = addPaths(paths)
	if err != nil {
		return err
	}

	return nil
}
