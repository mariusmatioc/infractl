package pkg

import (
	"fmt"
	"github.com/mariusmatioc/infractl/pkg/global"
	"path/filepath"
)

func Destroy(craftPath string, autoApprove bool) error {
	err := destroyCraft(craftPath, autoApprove)
	if err != nil {
		return err
	}
	fmt.Println("Destroy completed successfully")
	return nil
}

func destroyCraft(craftPath string, autoApprove bool) error {
	craft, err := global.NewCraft(craftPath)
	if err != nil {
		return err
	}
	destroy := []string{"destroy"}
	if autoApprove {
		destroy = append(destroy, "--auto-approve")
	}
	fmt.Println("Destroying ", craftPath)
	buildFolder, err := global.GetBuildFolder(global.NameOnly(craftPath))
	if err != nil {
		return err
	}
	err = global.RunTerraformCommand(buildFolder, destroy)
	if err != nil {
		return err
	}
	// Externals need to be destroyed last
	if ecs, ok := craft.(*global.EcsRecipe); ok {
		for _, external := range ecs.Externals {
			extPath := global.ToAbsPathBasedOn(filepath.Dir(craftPath), external.CraftFile)
			err = destroyCraft(extPath, autoApprove)
			if err != nil {
				return err
			}
		}
	}
	return err
}
