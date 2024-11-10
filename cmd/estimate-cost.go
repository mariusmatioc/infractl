package cmd

import (
	"github.com/mariusmatioc/infractl/pkg"
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/spf13/cobra"
)

func estimateCostCommand() *cobra.Command {
	var estimateCostCmd = &cobra.Command{
		Use:   "estimate-cost [[<root folder>] <craft-filename>]",
		Short: `Estimates the cost of the user's infrastructure`,
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
			err := pkg.EstimateCost(craftPath)
			return err
		},
	}
	return estimateCostCmd
}
