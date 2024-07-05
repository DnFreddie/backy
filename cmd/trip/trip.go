/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package trip

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TripCmd represents the trip command
var TripCmd = &cobra.Command{
	Use:   "trip",
	Short: "Trip-wire mini scans ur directory and checksum it requiers absolute path",
	Long: `
	Scans the directroy and gives u the schema holded in a 
	sqlite db and then produces any changes that were created 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("trip called")

	},
}

func init() {
	TripCmd.AddCommand(scanCmd)
	TripCmd.AddCommand(addCmd)
}
