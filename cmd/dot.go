/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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

		err:= dot.DotCommand()
		//var tmp  []dot.Dotfile
		//dot.CopyTemp(tmp)



		if err != nil {

			return
		}

		
	},
}

func init() {
	rootCmd.AddCommand(dotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
