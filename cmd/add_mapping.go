package cmd

import (
	"fmt"
	"log"
	"nuget-login-cli/nuget"
	"os"

	"github.com/spf13/cobra"
)

var addMappingCmd = &cobra.Command{
	Use:   "add-mapping [<name> | <url>] <pattern> [options]",
	Short: "Add source mapping",
	Long:  "Add source mapping to the specified config file, or default if not specified",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var configPath = nuget.GetNugetConfigPath(Target, Verbose)
		fmt.Printf("Using config file: %s\n", configPath)

		configExists := true
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configExists = false
		}

		name, err := nuget.GetNameForNugetSource(configPath, args[0], Verbose)
		if err != nil {
			log.Fatalf("Error determining name for source: %s", err)
		}

		if Defaults && configExists {
			fmt.Println("Warning - config already exists, default sources and mappings will not be added")
		}

		if Defaults && !configExists {
			fmt.Println("Adding default sources and mappings to new config file")
			nuget.AddSourceToNugetConfig(configPath, "nuget.org", "https://api.nuget.org/v3/index.json", Verbose)
			nuget.AddMappingToNugetConfig(configPath, "nuget.org", "*", Verbose)
		}

		fmt.Printf("Adding mapping for %s to %s..\n", name, args[1])
		err = nuget.AddMappingToNugetConfig(configPath, name, args[1], Verbose)
		if err != nil {
			log.Fatalf("Error adding mapping: %s", err)
		}

		fmt.Printf("Successfully added mapping for %s to %s\n", name, args[1])
	},
}
