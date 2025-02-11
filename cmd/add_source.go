package cmd

import (
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
		var configPath = nuget.GetNugetConfigPath(Target)

		var err error
		var name string
		var sourceUrl string
		if len(args) == 1 {
			sourceUrl = args[0]
			name, err = nuget.GetNameForNugetSource(configPath, sourceUrl)
			if err != nil {
				log.Fatalf("Error getting name for source: %s", err)
			}
		} else {
			name = args[0]
			sourceUrl = args[1]
		}

		nuget.AddSourceToNugetConfig(configPath, name, sourceUrl)

		if username != "" && password != "" {
			nuget.AddPackageSourceCredentialsToNugetConfig(configPath, name, username, password)
		}
	},
}

func init() {
	addSourceCmd.Flags().StringVarP(&username, "username", "u", "", "Username for the source")
	addSourceCmd.Flags().StringVarP(&password, "password", "p", "", "Password for the source")
	addSourceCmd.MarkFlagsRequiredTogether("username", "password")
}
