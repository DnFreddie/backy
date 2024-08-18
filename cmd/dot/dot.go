/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package dot

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
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
		fmt.Println("\nUre dots has been imported checkout them :)")
	},
}

func init() {
	DotCmd.Flags().StringVarP(&configPath, "path", "p", "", "specyfie the dotfiels target dir can be github url ")
	DotCmd.MarkFlagRequired("path")
	DotCmd.AddCommand(RevertCmd)
}

func dotCommand(repoPath string) error {
	isURL := isUrl(repoPath)
	id := uuid.New()
	r := &Repo{RepoId: fmt.Sprintln(id)}
	if isURL {
		err := r.Clone(repoPath)
		if err != nil {
			return fmt.Errorf("Failed to clone repo %s", r.RepoName)
		}
	} else {

		err := r.ReadLocal(repoPath)
		if err != nil {
			fmt.Println(err)
			return err
		}

	}
	err := r.getDots()
	if err != nil {
		return fmt.Errorf("error getting paths: %w", err)
	}
	force := true
	r.Link(force)
	
	r.saveRepoSchema()
	r.PrintRaport()
	return nil
}

func isUrl(str string) bool {

	if strings.Contains(str, "git@") {
		return true
	}
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
