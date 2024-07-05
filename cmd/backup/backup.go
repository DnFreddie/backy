package backup

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"time"

	"github.com/DnFreddie/backy/cmd/revert"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backuped bool
var archive bool
var BACKUP_DIR string
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Add the patsh that can be later backuped ",
	Long: `
	Add the paths that can be later backuped
	`,
	Run: func(cmd *cobra.Command, args []string) {
		Add_command(&args)
		if len(args) == 0 {
			cmd.Help()

		} else {

			if backuped {
				BACKUP_DIR= viper.GetViper().GetString("backup_dir")	
				Back(&args)

			}
			if archive {
				now := time.Now().Format("20060102150405")
				zipPath := fmt.Sprintf("%v.zip", now)
				err := ZipDir(args, zipPath)
				if err != nil {

					fmt.Println(err)
				}
			}
		}
	},
}

func init() {
	BackupCmd.Flags().BoolVarP(&backuped, "back", "b", false, "instant backup")
	BackupCmd.Flags().BoolVarP(&archive, "archive", "a", false, "archived the paths")
	BackupCmd.AddCommand(DeamonCmd)
	BackupCmd.AddCommand(revert.RevertCmd)

}
