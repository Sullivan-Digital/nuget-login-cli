package cmd

import (
	"fmt"
	"log"
	"nuget-login-cli/nuget"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new config file",
	Long:  "Initialize a new config file",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var configPath = nuget.GetNugetConfigPath(Target, Verbose)
		fmt.Printf("Initializing config file: %s\n", configPath)

		configExists := true
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configExists = false
		}

		if !configExists {
			err := nuget.InitializeEmptyNugetConfig(configPath, Verbose)
			if err != nil {
				log.Fatalf("Error initializing config: %s", err)
			}
		}

		if Defaults && configExists {
			fmt.Println("Warning - config already exists, default sources and mappings will not be added")
		}

		if Defaults && !configExists {
			fmt.Println("Adding default sources and mappings to new config file")
			err := nuget.AddSourceToNugetConfig(configPath, "nuget.org", "https://api.nuget.org/v3/index.json", Verbose)
			if err != nil {
				log.Fatalf("Error adding source: %s", err)
			}

			err = nuget.AddMappingToNugetConfig(configPath, "nuget.org", "*", Verbose)
			if err != nil {
				log.Fatalf("Error adding mapping: %s", err)
			}

			fmt.Printf("Successfully added default sources and mappings to %s\n", configPath)
		}
	},
}
