package add
/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Indexes the directory and adds it to the existing index group for future scanning",
	Long: `
	If there's no path specified for scanning, it adds the directory to the pool.
	It then checks each file, computes checksums, and stores them in the database for comparison during subsequent scans.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if addPath != "" {
			_, err := os.Stat(addPath)
			if err != nil {
				fmt.Println("the driectry doesn't exist")
				os.Exit(1)

			}
			TripAdd(addPath)
		}

	},
}

var (
	addPath string
)

func init() {
	AddCmd.Flags().StringVarP(&addPath, "path", "p", "", "the path to add to the scan")
	AddCmd.MarkFlagRequired("path")

}
