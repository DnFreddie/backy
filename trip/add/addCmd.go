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
