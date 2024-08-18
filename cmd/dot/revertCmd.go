package dot

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var options bool
var RevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Reverets back to the specyfied previous files cofiguration",
	Long: `Brings the confguration to the previouse state by deleting all the files
	and moving back the one that existed before.It also deletes the specyfied backup directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var backupToRevert string
		parent := cmd.Parent().Name()
		switch parent {
		case "dot":
			backupToRevert = viper.GetViper().GetString("config_dir")
		case "backup":
			backupToRevert = viper.GetViper().GetString("backup_dir")	
		default:
			log.Fatal("This shoud't have happedn this is the wrong parent", parent)

		}

		if cmd.Flags().Changed("delete") {
			//deletBackup(backupToRevert)
			fmt.Println(backupToRevert)

		} else {

			//r := &Repo{}

			dirs ,err := os.ReadDir("/home/rocky/.user_log/back_conf/")
			bV,err:= chooseBackupVersion(dirs)
			if err != nil  {
				log.Fatal("Failed to chose the veriosn ",err)
			}
			r,err := processReversion(path.Join("/home/rocky/.user_log/back_conf/",bV))

			if err != nil {
				
				log.Fatal(err)
			}
			fmt.Println(r)

			for _,i := range *r.Dots{
				fmt.Println(i)
			}




		}

	},
}

func init() {

	RevertCmd.Flags().BoolVarP(&options, "delete", "d", false, "Delete the chosen config backup")

}

