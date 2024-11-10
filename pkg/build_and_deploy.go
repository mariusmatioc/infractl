package pkg

import (
	"fmt"
	"github.com/mariusmatioc/infractl/pkg/global"
	"path/filepath"
)

// BuildAndDeploy builds and deploys all crafts starting with the given top-level craft
func BuildAndDeploy(craftPath string) (err error) {
	craft, err := global.NewCraft(craftPath)
	if err != nil {
		return
	}

	// The Terraform outputs to be used as environment variables in other crafts
	// The first key is the craft name, the second key is the output name
	outputs := make(map[string]map[string]string)

	if ecs, ok := craft.(*global.EcsRecipe); ok {
		// BuildAndDeploy externals first
		for _, external := range ecs.Externals {
			extPath := global.ToAbsPathBasedOn(filepath.Dir(craftPath), external.CraftFile)
			extRecipe, err2 := global.NewEcsCraft(extPath)
			if err2 != nil {
				return err2
			}
			err = extRecipe.BuildAndDeploy(extPath, external.Name, outputs)
			if err != nil {
				return
			}
		}
		// BuildAndDeploy the main craft
		err = ecs.BuildAndDeploy(craftPath, "", outputs)
		if err != nil {
			return
		}
		if len(outputs) > 0 {
			fmt.Println("Collected Terraform outputs:")
			for k, v := range outputs {
				fmt.Printf("External %s\n", k)
				for k2, v2 := range v {
					fmt.Printf("  %s=%s\n", k2, v2)
				}
			}
		}
		err = global.ListServices(craftPath)
	} else if net, ok := craft.(*global.NetworkRecipe); ok {
		// This is a network craft
		err = net.BuildAndDeploy(craftPath, outputs)
	} else if lam, ok := craft.(*global.LambdaRecipe); ok {
		// This is a lambda craft
		err = lam.BuildAndDeploy(craftPath, outputs)
	} else {
		err = fmt.Errorf("unknown craft type in %s", craftPath)
	}
	return
}
