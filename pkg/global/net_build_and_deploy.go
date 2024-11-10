package global

import (
	"fmt"
	"os"
)

// BuildAndDeploy builds the Terraform files, deploys them and collects outputs into the outputs map
func (net *NetworkRecipe) BuildAndDeploy(craftPath string, outputs map[string]map[string]string) error {
	craftName := NameOnly(craftPath)
	// Folder where the Terraform files will be generated
	buildFolder, err := GetBuildFolder(craftName)
	if err != nil {
		return err
	}

	fmt.Printf("Writing Terraform files to '%s'...\n", buildFolder)
	err = DeleteFiles(buildFolder, "*.tf")
	if err != nil {
		return err
	}

	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	err = os.Chdir(RootParent) // we want relative file names to refer to parent of root folder
	if err != nil {
		return err
	}

	conf := Config{}
	conf.Recipe = net
	conf.BuildFolder = buildFolder
	conf.OutputsMap = outputs

	err = conf.BuildNetwork()
	if err != nil {
		return err
	}

	//net.SetAwsEnvs()
	err = TerraformDeploy(buildFolder, craftName)
	if err != nil {
		return err
	}

	// Collect outputs
	outs := make(map[string]string)
	err = GetTerraformOutputs(buildFolder, outs)
	if err != nil {
		return err
	}
	fmt.Printf("'%s' had %d outputs\n", craftName, len(outs))
	if len(outs) > 0 {
		outputs[craftName] = outs
	}
	return nil
}
