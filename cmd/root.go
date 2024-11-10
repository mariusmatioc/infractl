// Package cmd exposes the CLI commands
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  exeName + " - a CLI for cloud deployment",
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	rootCmd.AddCommand(versionCommand())
	rootCmd.AddCommand(initCommand())
	//rootCmd.AddCommand(remoteStateCommand())
	rootCmd.AddCommand(deployCommand())
	rootCmd.AddCommand(destroyCommand())
	rootCmd.AddCommand(outputsCommand())
	rootCmd.AddCommand(listClusters())
	rootCmd.AddCommand(listServices())
	rootCmd.AddCommand(estimateCostCommand())
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}
