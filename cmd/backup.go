/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	//"github.com/DnFreddie/backy/backup"
	"github.com/DnFreddie/backy/backup"
	"github.com/DnFreddie/backy/backup/deamon"
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
		backup.Add_command(&args)
		if len(args) == 0 {

			cmd.Help()
		} else {

			if backuped {
				backup.Back(&args)

			}
			if archive {
				now := time.Now().Format("20060102150405")
				zipPath:=fmt.Sprintf("%v.zip",now)
				err := backup.ZipDir(args, zipPath)
				if err != nil {

					fmt.Println(err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().BoolVarP(&backuped, "back", "b", false, "instant backup")
	backupCmd.Flags().BoolVarP(&archive, "archive", "a", false, "archived the paths")
	backupCmd.AddCommand(deamon.DeamonCmd)
	
}
