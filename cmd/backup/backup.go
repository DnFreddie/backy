package backup

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backuped bool
var archive bool
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
				backup_dir:= viper.GetViper().GetString("backup_dir")	
				Back(&args,backup_dir)

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

}
