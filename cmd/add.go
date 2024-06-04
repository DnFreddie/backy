/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/DnFreddie/backy/trip"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if addPath != "" {
			_, err := os.Stat(addPath)
			if err != nil {
				fmt.Println("the driectry doesn't exist")
				os.Exit(1)

			}
			trip.TripAdd(addPath)
		}

	},
}

var (
	addPath string
)

func init() {
	tripCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addPath, "path", "p", "", "the path to add to the scan")
	addCmd.MarkFlagRequired("path")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
