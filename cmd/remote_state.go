package cmd

import (
	"fmt"
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/mariusmatioc/infractl/pkg"
	"github.com/spf13/cobra"
)

func remoteStateCommand() *cobra.Command {
	var remoteStateCmd = &cobra.Command{
		Use:   "remote-state [root folder] create|destroy <project-name>",
		Short: fmt.Sprintf("Creates or destroys the remote state for the project"),
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsIx := 0
			if len(args) == 3 {
				global.SetRootFolder(args)
				argsIx = 1
			} else {
				global.SetDefaultRootFolder()
			}
			create := true
			switch args[argsIx] {
			case "create":
				create = true
			case "destroy":
				create = false
			default:
				return fmt.Errorf("invalid action %s. Must be create or destroy", args[argsIx])
			}
			return pkg.RemoteState(create, args[argsIx+1])
		},
	}
	return remoteStateCmd
}
