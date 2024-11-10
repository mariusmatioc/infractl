package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func versionCommand() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(VERSION)
			return nil
		},
	}
	return versionCmd
}
