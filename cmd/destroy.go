package cmd

import (
	"github.com/mariusmatioc/infractl/pkg"
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/spf13/cobra"
)

func destroyCommand() *cobra.Command {
	var destroyCmd = &cobra.Command{
		Use:   "destroy [<root folder>] <craft-filename>",
		Short: `Destroys previously created infrastructure on AWS`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			craftPath, err := global.GetCraftPath(args)
			if err == nil {
				autoApprove, err2 := cmd.Flags().GetBool(AUTO_APPROVE)
				if err2 != nil {
					return err2
				}
				err = pkg.Destroy(craftPath, autoApprove)
			}
			return err
		},
	}
	destroyCmd.PersistentFlags().BoolP(AUTO_APPROVE, "a", false, "")
	return destroyCmd
}
