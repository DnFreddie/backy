package dot

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/

import (
	"fmt"
	"os"

	"github.com/DnFreddie/backy/common"
	"github.com/DnFreddie/backy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dot/deleteCmd represents the dot/delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")
		DOTS = viper.GetViper().GetString("dots")
		dots, err := utils.Checkdir(DOTS, false)
		fmt.Println("Thsi are DOTS", DOTS)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		r := &Repo{}

		common.DeleteBackup(r, dots)

	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dot/deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dot/deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
