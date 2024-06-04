/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// tripCmd represents the trip command
var tripCmd = &cobra.Command{
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
	rootCmd.AddCommand(tripCmd)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tripCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
