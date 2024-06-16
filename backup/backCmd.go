package backup

import "fmt"

func Add_command(args *[]string) error {
	paths, err := Add_dir(args)

	if err != nil {
		return err
	}
	if len(paths)==0{
		fmt.Println("Skipping... Nothing to add ")
		return nil
	}
	err = Jsonyfie(paths)
	if err != nil {
		return err
	}

	return nil
}
