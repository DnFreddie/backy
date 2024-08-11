/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package dot

import (
	"fmt"
	"github.com/DnFreddie/backy/cmd/revert"
	"github.com/DnFreddie/backy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
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
		fmt.Println("\nUre dots has been imported checkout them :)")
	},
}

func init() {
	DotCmd.Flags().StringVarP(&configPath, "path", "p", "", "specyfie the dotfiels target dir can be github url ")
	DotCmd.MarkFlagRequired("path")
	DotCmd.AddCommand(revert.RevertCmd)
}

func dotCommand(local string) error {
	var dest string
	dest = local
	isURL := isUrl(local)

	r := &Repo{}
	if isURL {
		err := r.Clone(local)
		if err != nil{
			fmt.Errorf("Failed to clone repo ",r.RepoName)
		}
		dest = r.AbsP
	} 

	absDest, err := utils.MakeAbsolute(dest)
	if err != nil {
		return fmt.Errorf("%v doesn't exist: %w", path.Base(dest), err)
	}

	dotStructs, err := getPaths(absDest)
	if err != nil {
		return fmt.Errorf("error getting paths: %w", err)
	}

	for _, dot := range *dotStructs {
		dot.IsExe()
	}

	if err := createSymlink(*dotStructs, absDest); err != nil {
		return fmt.Errorf("error creating symlink: %w", err)
	}

	return nil
}

func getPaths(gitPath string) (*[]Dotfile, error) {

	dirs, err := os.ReadDir(gitPath)
	if err != nil {
		fmt.Println("Can't list this dir probably permissions issue ", err)
		return nil, err

	}
	var dotfiels []Dotfile

	for _, d := range dirs {

		dot := Dotfile{
			Location: d,
			AbPath: path.Join(gitPath,d.Name()),

		}
		dotfiels = append(dotfiels, dot)

	}

	toIgnore := readIgnore()
	for _, dot := range dotfiels {
		dot.ignore(&toIgnore)
	}

	return &dotfiels, nil
}
