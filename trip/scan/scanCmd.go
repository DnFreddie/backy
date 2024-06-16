/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scan

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := TripScan(csvName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
var (
	csvName string
)

func init() {
	ScanCmd.Flags().StringVarP(&csvName, "csv", "c", "", "wirte to a given csv path")
}
