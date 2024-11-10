package pkg

import (
	"fmt"
	"github.com/mariusmatioc/infractl/pkg/global"
)

func Outputs(craftPath string) error {
	craftNames := make(map[string]bool)
	err := global.CollectCraftNames(craftPath, craftNames)
	if err != nil {
		return err
	}
	for craftName := range craftNames {
		path, err := global.GetAbsoluteCraftPath(craftName)
		if err != nil {
			return err
		}
		err = outputsForCraft(path)
		if err != nil {
			return err
		}
	}
	return err

}

func outputsForCraft(craftPath string) error {
	buildFolder, err := global.GetBuildFolder(global.NameOnly(craftPath))
	if err != nil {
		return err
	}
	fmt.Println("Outputs for", craftPath)
	err = global.RunTerraformCommand(buildFolder, []string{"output"})
	fmt.Println()
	if err != nil {
		return err
	}
	return err
}
