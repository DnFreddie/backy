/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/DnFreddie/backy/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backy",
	Short: "Backup, Dotfiles and Tripwire all in one binary",
	Long: `
This application combines the functionality of backup management, dotfiles synchronization, and file integrity monitoring into a single versatile binary.

1. Backup Management:
   Easily configure and initiate backups of important files and directories.
	Customize backup schedules and destinations using a simple configuration file.

2. Dotfiles Synchronization:
   Streamline the management of your dotfiles across multiple systems.
	Ensure consistency and synchronization of configurations for applications like shells, editors, and more.

3. File Integrity Monitoring (Tripwire):
   Monitor changes in critical files and directories using cryptographic hash comparisons.
	Detect unauthorized modifications and maintain the integrity of your system.
For more details and usage examples, refer to the documentation or run the application with '--help'.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.LoadConfig)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
