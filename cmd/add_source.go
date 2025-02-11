package cmd

import (
	"fmt"
	"log"
	"nuget-login-cli/nuget"

	"github.com/spf13/cobra"
)

var (
	username string
	password string
)

var addSourceCmd = &cobra.Command{
	Use:   "add-source [<url> | <name> <url>] [options]",
	Short: "Add nuget source",
	Long:  "Add nuget source to the specified config file, or default if not specified",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var configPath = nuget.GetNugetConfigPath(Target, Verbose)
		fmt.Printf("Using config file: %s\n", configPath)

		var err error
		var name string
		var sourceUrl string
		if len(args) == 1 {
			sourceUrl = args[0]
			name, err = nuget.GetNameForNugetSource(configPath, sourceUrl, Verbose)
			if err != nil {
				log.Fatalf("Error getting name for source: %s", err)
			}
		} else {
			name = args[0]
			sourceUrl = args[1]
		}

		fmt.Printf("Adding source %s to %s..\n", name, sourceUrl)
		err = nuget.AddSourceToNugetConfig(configPath, name, sourceUrl, Verbose)
		if err != nil {
			log.Fatalf("Error adding source: %s", err)
		}

		if username != "" && password != "" {
			fmt.Printf("Adding package source credentials for %s..\n", name)
			err = nuget.AddPackageSourceCredentialsToNugetConfig(configPath, name, username, password, Verbose)
			if err != nil {
				log.Fatalf("Error adding package source credentials: %s", err)
			}

			fmt.Printf("Successfully added source and credentials for %s to %s\n", name, sourceUrl)
		} else {
			fmt.Printf("Successfully added source %s to %s\n", name, sourceUrl)
		}
	},
}

func init() {
	addSourceCmd.Flags().StringVarP(&username, "username", "u", "", "Username for the source")
	addSourceCmd.Flags().StringVarP(&password, "password", "p", "", "Password for the source")
	addSourceCmd.MarkFlagsRequiredTogether("username", "password")
}
