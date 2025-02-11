package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nuget-login-cli",
	Short: "nuget-login-cli is a CLI tool to login to a NuGet registry",
}

var (
	Target string
	Verbose bool
	Defaults bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


func init() {
	rootCmd.PersistentFlags().StringVarP(&Target, "target", "t", "", "target config file (default is $HOME/.nuget/nuget.config)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Defaults, "with-defaults", "", false, "add default sources and mappings (nuget.org) to new config files. Does not affect existing config files.")

	rootCmd.AddCommand(addSourceCmd)
	rootCmd.AddCommand(addMappingCmd)
	rootCmd.AddCommand(initCmd)
}