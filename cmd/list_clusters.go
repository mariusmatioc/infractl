package cmd

import (
	"github.com/mariusmatioc/infractl/pkg"
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/spf13/cobra"
)

func listClusters() *cobra.Command {
	var listClustersCmd = &cobra.Command{
		Use:   "list-clusters [<root folder>] [<craft-filename>]",
		Short: `Lists the user's ECS clusters`,
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			craftPath := ""
			if len(args) > 0 {
				var err error
				craftPath, err = global.GetCraftPath(args)
				if err != nil {
					return err
				}
			}
			err := pkg.ListClusters(craftPath)
			return err
		},
	}
	return listClustersCmd
}
