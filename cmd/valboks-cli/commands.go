package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"valboks/pkg/dropbox"
)

func newAuthCommand() *cobra.Command {
	var appKey, appSecret string

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticat with Dropbox",
		Long:  "TBA",
		RunE: func(cmd *cobra.Command, args []string) error {
			printVerbose(cmd, "Starting authentication process")

			if appKey == "" || appSecret == "" {
				return fmt.Errorf("Starting authentication process")
			}

			fmt.Println("📋 Dropbox Authentication")
			fmt.Println("=========================")
			fmt.Printf("App Key: %s\n", appKey)
			fmt.Println()
			fmt.Println("To complete authentication:")
			fmt.Println("1. Visit https://www.dropbox.com/developers/apps")
			fmt.Println("2. Find your app and go to the 'Settings' tab")
			fmt.Println("3. Generate an access token")
			fmt.Println("4. Enter the access token below")
			fmt.Println()

			fmt.Print("Enter your access token: ")
			var accessToken string
			_, err := fmt.Scanln(&accessToken)
			if err != nil {
				return fmt.Errorf("error reading access token: %w", err)
			}

			fmt.Printf("Access token: %s\n", accessToken)

			accessToken = strings.TrimSpace(accessToken)
			if accessToken == "" {
				return fmt.Errorf("access token cannot be empty")
			}

			printVerbose(cmd, "Testing connection with provided token")

			client := dropbox.NewClient(accessToken)
			err = client.TestConnection()
			if err != nil {
				return fmt.Errorf("authentication to dropbox failed - invalid token: %w", err)
			}

			//Saving the creds
			configManager.SetCredentials(appKey, appSecret, accessToken)
			err = configManager.Save()
			if err != nil {
				return fmt.Errorf("error saving configuration: %w", err)
			}

			fmt.Println("✅ Authentication successful!")
			printVerbose(cmd, "Configuration saved successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&appKey, "app-key", "", "Dropbox API app key (required)")
	cmd.Flags().StringVar(&appSecret, "app-secret", "", "Dropbox API app secret (required)")
	cmd.MarkFlagRequired("app-key")
	cmd.MarkFlagRequired("app-secret")

	return cmd
}

func newListCommand() *cobra.Command {

	var longFormat bool

	cmd := &cobra.Command{
		Use:     "ls [path]",
		Aliases: []string{"list"},
		Short:   "List files and folders",
		Long:    `List files and folders in the specified Dropbox path.`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !configManager.IsConfigured() {
				return fmt.Errorf("not authenticated - run 'auth' command first")
			}

			path := "/"
			if len(args) > 0 {
				path = args[0]
			}

			printVerbose(cmd, "Listing contents of: %s", path)

			client := dropbox.NewClient(configManager.GetConfig().AccessToken)
			fileInfos, err := client.ListFolder(path)
			if err != nil {
				return err
			}

			if len(fileInfos) == 0 {
				fmt.Println("📂 Empty folder")
				return nil
			}

			printVerbose(cmd, "Found %d items", len(fileInfos))

			for _, info := range fileInfos {
				if longFormat {
					if info.IsFolder {
						fmt.Printf("📁 %-30s <DIR>\n", info.Name)
					} else {
						fmt.Printf("📁 %-30s %d bytes\n", info.Name, info.Size)
					}
				} else {
					if info.IsFolder {
						fmt.Printf("📁 %s\n", info.Name)
					} else {
						fmt.Printf("📁 %s\n", info.Name)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&longFormat, "long", "l", false, "Use long listing format")

	return cmd
}

func newUploadCommand() *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:     "put [local_path] [dropbox_path]",
		Aliases: []string{"upload"},
		Short:   "Upload a file to Dropbox",
		Long:    `Upload a file from your local filesystem to Dropbox.`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !configManager.IsConfigured() {
				return fmt.Errorf("not authenticated - run 'auth' command first")
			}

			localPath := args[0]
			dropboxPath := args[1]

			//Check if local file exists
			if _, err := os.Stat(localPath); os.IsNotExist(err) {
				return fmt.Errorf("local file '%s' does not exist", localPath)
			}

			printVerbose(cmd, "Uploading %s to %s (overwrite: %v)", localPath, dropboxPath, overwrite)

			client := dropbox.NewClient(configManager.GetConfig().AccessToken)
			err := client.UploadFile(localPath, dropboxPath, overwrite)
			if err != nil {
				fmt.Errorf("Failed to up load the file to dropbox")
				return err
			}

			fmt.Printf("✅ Uploaded '%s' to '%s'\n", localPath, dropboxPath)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")

	return cmd
}

func newDeleteCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "rm [path]",
		Aliases: []string{"delete"},
		Short:   "Delete a file or folder",
		Long:    `Delete a file or folder from Dropbox.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !configManager.IsConfigured() {
				return fmt.Errorf("not authenticated - run 'auth' command first")
			}

			path := args[0]

			if !force {
				fmt.Printf("Are you sure you want to delete '%s'?, (Y/N): ", path)
				var response string

				fmt.Scanln(&response)
				if strings.ToLower(response) != "Y" && strings.ToLower(response) != "yes" {
					fmt.Println("Deletion cancelled")
					return nil
				}
			}

			printVerbose(cmd, "Deleting: %s", path)

			client := dropbox.NewClient(configManager.GetConfig().AccessToken)
			err := client.DeletePath(path)
			if err != nil {
				fmt.Errorf("failed to delete the file")
				return nil
			}

			fmt.Printf("✅ Deleted '%s'\n", path)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")

	return cmd
}

func newInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [path]",
		Short: "Get information about a file or folder",
		Long:  `Get detailed information about a file or folder in Dropbox,`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !configManager.IsConfigured() {
				return fmt.Errorf("not authenticated - run 'auth' command first")
			}

			path := args[0]

			printVerbose(cmd, "Getting info for: %s", path)

			client := dropbox.NewClient(configManager.GetConfig().AccessToken)
			info, err := client.GetFileInfo(path)
			if err != nil {
				fmt.Errorf("failed to get the file Info: %+v", err)
				return err
			}

			fmt.Printf("📋 Information for '%s'\n", path)
			fmt.Printf("	Name: %s\n", info.Name)
			fmt.Printf("	Path: %s\n", info.Path)
			if info.IsFolder {
				fmt.Printf("	Type: Folder\n")
			} else {
				fmt.Printf("	Type: File\n")
				fmt.Printf("	Size: %d bytes\n", info.Size)
			}

			return nil
		},
	}

	return cmd
}
