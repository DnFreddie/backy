/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package scan

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scanning the directory and checking the difference in the hash of the files",
	Long: `

Tripwire main functionality:
It scans the specified directory to detect any changes in the files' content.
The tool then outputs the current state in CSV format.

`,
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
