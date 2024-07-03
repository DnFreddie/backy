/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package deamon

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deamonCmd represents the deamon command
var DeamonCmd = &cobra.Command{
	Use:   "deamon",
	Short: "Start a backy  deamon that backups the files daily ",
	Long: `
		Starts a backy daemon to backup files as they are added.
		By default, backups occur daily; you can modify this setting in .backy.yaml using the cron_time field.`,
	
	Run: func(cmd *cobra.Command, args []string) {
		c:= viper.GetViper().GetString("cron_time")
		fmt.Println("Deamon started ")
		loopEternally(c)
	},
}

func init() {
	//.AddCommand(deamonCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deamonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deamonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
