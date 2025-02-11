package cmd

import (
	"log"
	"nuget-login-cli/nuget"

	"github.com/spf13/cobra"
)

var addMappingCmd = &cobra.Command{
	Use:   "add-mapping [<name> | <url>] <pattern> [options]",
	Short: "Add source mapping",
	Long:  "Add source mapping to the specified config file, or default if not specified",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var configPath = nuget.GetNugetConfigPath(Target)

		name, err := nuget.GetNameForNugetSource(configPath, args[0])
		if err != nil {
			log.Fatalf("Error determining name for source: %s", err)
		}

		nuget.AddMappingToNugetConfig(configPath, name, args[1])
	},
}

func init() {

}
