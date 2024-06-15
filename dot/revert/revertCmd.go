package revert

import (
	"log"

	"github.com/spf13/cobra"
)

var options bool
var RevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "revert",
	Long:  "hahahh",
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

	RevertCmd.Flags().BoolVarP(&options, "date", "d", false, "Enable debug mode")

}
