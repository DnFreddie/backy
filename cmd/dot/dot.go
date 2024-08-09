/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package dot

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"github.com/DnFreddie/backy/cmd/revert"
	"github.com/DnFreddie/backy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var BACK_CONF string
var TARGET string
var configPath string

const (
	IGNORE = ".gitignore"
)

var DotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Imports docs files or intsall specyifed",
	Long: `
	Import the docs fiels from the sepcyfied repo and or intall 
	them from any other  directory.
	The repo can be specyfied either by the .dotfiels fiel or direcly in command 
	It'cheks for the binary programs in ure user paths and based on that imports 
	x
	f the already existing files and creates symlinks to .config 
	It can also isntall the one that not exist using -i 

	`,
	Run: func(cmd *cobra.Command, args []string) {

		//Where the app is looking  default .cofig
		TARGET = viper.GetViper().GetString("config_path")
		// Backup dir for the configs
		BACK_CONF = viper.GetViper().GetString("config_dir")

		err := dotCommand(configPath)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Ure dots has been imported checkout them :)")
	},
}

func init() {
	DotCmd.Flags().StringVarP(&configPath, "path", "p", "", "specyfie the dotfiels target dir can be github url ")
	DotCmd.MarkFlagRequired("path")
	DotCmd.AddCommand(revert.RevertCmd)
}

func dotCommand(repo string) error {
	var URL bool
	var dest string
	URL = isUrl(repo)
	if URL {
		clonedDest, err := gitClone(repo)
		if err != nil {
			log.Fatal("Failed to copy url")
		}
		dest = clonedDest
	} else {
		dest = repo
	}

	absDest, err := utils.MakeAbsoulute(dest)

	if err != nil {
		log.Fatalf("%v doesn't exist\n", path.Base(dest))
	}

	dotStructs, err := getPaths(absDest)
	if err != nil {
		fmt.Println(err)
		return err
	}


	for _, dot := range dotStructs { dot.IsExe() }


	err = createSymlink(dotStructs, absDest)
	if err != nil {
		return err
	}

	return nil
}

func getPaths(gitPath string) ([]Dotfile, error) {
	dirs, err := os.ReadDir(gitPath)
	if err != nil {
		fmt.Println("Can't list this dir probably permissions issue ", err)
		return nil, err

	}
	var dotfiels []Dotfile

	for _, d := range dirs {

		dot := Dotfile{
			Location: d,
			Repo:     "test",
		}
		dotfiels = append(dotfiels, dot)

	}

	toIgnore, err := readIgnore()
	if err != nil {
		fmt.Println("Can't read git ignore: ", err)
		return nil, err
	}

	var paths []fs.DirEntry
	for _, dir := range dirs {
		if !shouldIgnore(dir.Name(), toIgnore) {
			paths = append(paths, dir)
		}
	}
	
	return dotfiels, nil
}
