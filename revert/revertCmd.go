package revert

import (
	"log"

	"github.com/DnFreddie/backy/backup"
	"github.com/DnFreddie/backy/dot"
	"github.com/spf13/cobra"
)

var options bool
var RevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Reverets back to the specyfied previous files cofiguration",
	Long:  `Brings the confguration to the previouse state by deleting all the files
	and moving back the one that existed before.It also deletes the specyfied backup directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var backupToRevert string
		parent:= cmd.Parent().Name()
		switch parent{
		case "dot":
		backupToRevert=dot.BACK_CONF
		case "backup":
		backupToRevert=backup.BACKUP_PATH
		default:

		log.Fatal("This shoud't have happedn this is the wrong parent",parent)

		}

		if cmd.Flags().Changed("delete") {
			deletBackup(backupToRevert)

		} else {

			err := revertBackups(backupToRevert)
			if err != nil {
				log.Fatal(err)
			}

		}

	},
}

func init() {

	RevertCmd.Flags().BoolVarP(&options, "delete", "d", false, "Delete the chosen config backup")

}
