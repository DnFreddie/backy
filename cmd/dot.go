/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/DnFreddie/backy/dot"
	"github.com/spf13/cobra"
)

// dotCmd represents the dot command
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


		//err := dot.RevertBackups()
		err:= dot.DotCommand(addPath)
		//fmt.Println("the changes has been suscefully reverted ")
		fmt.Println("Ure dots has been imported checkout them :)")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//var tmp  []dot.Dotfile
		//dot.CopyTemp(tmp)

	},
}

var revertCmd = &cobra.Command{
		Use:   "revert",
		Short: "revert",
		Long:  "hahahh",
		Run: func (cmd *cobra.Command, args []string) {
			// Command implementation here
		dot.RevertBackups()

		},
	}







func init() {
	rootCmd.AddCommand(dotCmd)
	dotCmd.Flags().StringVarP(&addPath, "path", "p", "", "specyfie the dotfiels target dir can be github url ")
	dotCmd.AddCommand(revertCmd)
	dotCmd.MarkFlagRequired("path")
	// is called directly, e.g.:
	// dotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
