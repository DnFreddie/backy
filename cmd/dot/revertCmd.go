package dot

import (
	"log"
	"github.com/DnFreddie/backy/common"
	"github.com/DnFreddie/backy/utils"
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
		DOTS = viper.GetViper().GetString("dots")
		backupDir, err := utils.Checkdir(DOTS, false)
		if err != nil {
			log.Fatal(utils.UserError{FPath: DOTS, Err: err}.Err.Error())
		}
		r := &Repo{}
		err = common.Revert(r, backupDir)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {

}
