package cmd

import (
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/spf13/cobra"
)

func listServices() *cobra.Command {
	var listServicesCmd = &cobra.Command{
		Use:   "list-services [<root folder>] <craft-filename>",
		Short: `Lists the services`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			craftPath, err := global.GetCraftPath(args)
			if err == nil {
				err = global.ListServices(craftPath)
			}
			return err
		},
	}
	return listServicesCmd
}
