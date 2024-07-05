/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package backup

import (
	"fmt"

	"github.com/DnFreddie/backy/utils"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var DeamonCmd = &cobra.Command{
	Use:   "deamon",
	Short: "Start a backy  deamon that backups the files daily ",
	Long: `
		Starts a backy daemon to backup files as they are added.
		By default, backups occur daily; you can modify this setting in .backy.yaml using the cron_time field.`,

	Run: func(cmd *cobra.Command, args []string) {
		c := viper.GetViper().GetString("cron_time")
		fmt.Println("Deamon started ")
		loopEternally(c)
	},
}

func init() {
}

func loopEternally(croneRecord string) {
	db, err := utils.InitDb(BACKUP_DB, Brecord{})
	if err != nil {
		log.Fatal("Failed to connect to the db:", err)
	}

	c := cron.New()

	_, err = c.AddFunc(croneRecord, func() {
		var records []Brecord
		var paths []string

		result := db.Find(&records)
		if result.Error != nil {
			log.Println("Failed to retrieve records:", result.Error)
			return
		}

		for _, record := range records {
			paths = append(paths, record.TargetPath)

		}
		fmt.Println(paths)
		err := Back(&paths)
		if err != nil {
			fmt.Println(err)
			return
		}

	})
	if err != nil {
		log.Fatal("Failed to add cron job:", err)
	}

	c.Start()

	select {}
}

