package cmd

import (
	"fmt"

	"github.com/mariusmatioc/infractl/pkg"
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init [folder]",
		Short: fmt.Sprintf(`Initializes the "%s" directory`, global.InfraCtlRoot),
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			global.SetRootFolder(args)
			return pkg.InitInfractl()
		},
	}
	return initCmd
}
