 package main

 import (
	 "fmt"
	 "github.com/spf13/cobra"
	 "strings"
	 "valboks/pkg/dropbox"
 )

 func newAuthCommand() *cobra.Command {
	var appKey, appSecret string

	cmd := &cobra.Command{
		Use: "auth",
		Short: "Authenticat with Dropbox",
		Long: "TBA",
		RunE: func(cmd * cobra.Command, args []string) error {
			printVerbose(cmd, "Starting authentication process")

			if appKey == "" || appSecret == "" {
				return fmt.Errorf("Starting authentication process")
			}

			fmt.Println("📋 Dropbox Authentication")
			fmt.Println("=========================")
			fmt.Printf("App Key: %s\n", appKey)
			fmt.Println()
			fmt.Println("To complete authetication:")
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
			printVerbose(cmd, "Configuration saved succefully")
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
		Use: "ls [path]",
		Aliases: []string{"list"},
		Short: "List files and folders",
		Long:  `List files and folders in the specified Dropbox path.`,
		Args: cobra.MaximumNArgs(1),
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