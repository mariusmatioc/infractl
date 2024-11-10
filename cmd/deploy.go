package cmd

import (
	"github.com/mariusmatioc/infractl/pkg"
	"github.com/mariusmatioc/infractl/pkg/global"
	"github.com/spf13/cobra"
)

// Flags
const (
	FORCE_REBUILD = "force-rebuild"
)

func deployCommand() *cobra.Command {
	var deployCmd = &cobra.Command{
		Use:   "deploy [root folder] <craft-filename>",
		Short: `Creates the Terraform files and deploys`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			craftPath, err := global.GetCraftPath(args)
			if err != nil {
				return err
			}
			global.ForceRebuild, _ = cmd.Flags().GetBool(FORCE_REBUILD)
			backendS, _ := cmd.Flags().GetString(BACKEND)
			if backendS != "" {
				err = global.SetBackend(backendS)
			}
			if err == nil {
				err = pkg.BuildAndDeploy(craftPath)
			}
			return err
		},
	}
	deployCmd.PersistentFlags().BoolP(FORCE_REBUILD, "f", false, "force new Docker image build")
	deployCmd.PersistentFlags().StringP(BACKEND, "b", "", "S3 backend id")
	return deployCmd
}
