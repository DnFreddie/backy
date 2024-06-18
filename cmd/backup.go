/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	//"github.com/DnFreddie/backy/backup"
	"github.com/DnFreddie/backy/backup"
	"github.com/spf13/cobra"
)

var backuped bool
var archive bool
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Add the patsh that can be later backuped ",
	Long: `
	Add the paths that can be later backuped
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {

		// If this is a file, write its contents to the zip archive
			cmd.Help()
		} else {

			if backuped {
				backup.Add_command(&args)
				fmt.Println("Backuped called")
				backup.Back(&args)

			}
			if archive{
				fmt.Println("archive")
				err := backup.ZipDir("/home/aura/Documents/Notes/","./ans.zip")
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().BoolVarP(&backuped, "back", "b", false, "instant backup")
	backupCmd.Flags().BoolVarP(&archive, "archive", "a", false, "archived the paths")

}
