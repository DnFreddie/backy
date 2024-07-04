/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/DnFreddie/backy/dot"
	"github.com/DnFreddie/backy/revert"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configPath string

var dotCmd = &cobra.Command{
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

		dot.TARGET = viper.GetViper().GetString("config_path")
		err := dot.DotCommand(configPath)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Ure dots has been imported checkout them :)")
	},
}

func init() {
	rootCmd.AddCommand(dotCmd)
	dotCmd.Flags().StringVarP(&configPath, "path", "p", "", "specyfie the dotfiels target dir can be github url ")
	dotCmd.AddCommand(revert.RevertCmd)

	dotCmd.MarkFlagRequired("path")
}
