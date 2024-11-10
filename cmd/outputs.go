package cmd

import (
	"github.com/mariusmatioc/infractl/pkg"
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/spf13/cobra"
)

func outputsCommand() *cobra.Command {
	var outputsCmd = &cobra.Command{
		Use:   "outputs [<root folder>] <craft-filename>",
		Short: `Shows the outputs of the Terraform deployment`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			craftPath, err := global.GetCraftPath(args)
			if err != nil {
				return err
			}
			err = pkg.Outputs(craftPath)
			return err
		},
	}
	return outputsCmd
}
