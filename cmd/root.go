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
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


func init() {
	rootCmd.PersistentFlags().StringVarP(&Target, "target", "t", "", "target config file (default is $HOME/.nuget/nuget.config)")

	rootCmd.AddCommand(addSourceCmd)
	rootCmd.AddCommand(addMappingCmd)
}