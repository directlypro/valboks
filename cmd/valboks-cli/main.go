package main

import (
	"fmt"
	"os"
	"valboks/internal/config"
	"github.com/spf13/cobra"
)

var (
	configManager *config.ConfigManager
	version = "1.0.0"
	commit = "dev"
	date = "2025-06-20"
)

func main() {
	var err error

	configManager, err = config.NewConfigManager()
	if err != nil {
		fmt.Printf("Error initializing configuration: %v\n", err)
		os.Exit(1)
	}

	err = configManager.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use: "valboks-cli",
		Short: "Custom Dropbox CLI tool",
		Long: "TBA",
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(newAuthCommand())
	rootCmd.AddCommand(newListCommand())
//	rootCmd.AddCommand(newDownloadCommand()) // Due to vibe coding this is not complete
	rootCmd.AddCommand(newUploadCommand())
	rootCmd.AddCommand(newDeleteCommand())
//	rootCmd.AddCommand(newMkdirCommand()) // Due to vibe coding this is not complete
	rootCmd.AddCommand(newInfoCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n, err")
		os.Exit(1)
	}
}

func getVerbose(cmd *cobra.Command) bool {
	verbose, _ := cmd.Flags().GetBool("verbose")
	return verbose
}

func printVerbose(cmd *cobra.Command, format string, args ...interface{}) {
	if getVerbose(cmd) {
		fmt.Printf("[VERBOSE] "+format+"\n", args...)
	}
}