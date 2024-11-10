package pkg

import (
	"github.com/mariusmatioc/infractl/pkg/global"
	"github.com/mariusmatioc/infractl/pkg/remote_state"
)

// RemoteState creates or destroys the remote state for the project
// This is optional, as one can use local state as well
func RemoteState(create bool, projectName string) error {
	buildFolder, err := global.GetBuildFolder(projectName + "-remote-state")
	if err != nil {
		return err
	}
	if create {
		err = global.DeleteFiles(buildFolder, "*.tf")
		if err != nil {
			return err
		}
		config := struct{ ProjectName string }{projectName}
		err = global.BuildFromTemplate(remote_state.Backend, buildFolder, "backend.tf", &config)
		if err != nil {
			return err
		}
		err = global.TerraformDeploy(buildFolder, projectName)
	} else {
		err = global.RunTerraformCommand(buildFolder, []string{"destroy"})
	}
	return err
}
