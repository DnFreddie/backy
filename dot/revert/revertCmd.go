package revert

import (
	"github.com/spf13/cobra"
	"log"
)

var options bool
var RevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Reverets back to the specyfied previous .config cofiguration",
	Long:  `Brings the confguration to the previouse state by deleting all the files
	and moving back the one that existed before.It also deletes the specyfied backup directory`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("date") {
			deletBackup()

		} else {

			err := RevertBackups()
			if err != nil {
				log.Fatal(err)
			}

		}

	},
}

func init() {

	RevertCmd.Flags().BoolVarP(&options, "delte", "d", false, "Delete the chosen config backup")

}
