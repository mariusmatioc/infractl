package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mariusmatioc/infractl/data/recipes"
	"github.com/mariusmatioc/infractl/pkg/global"
)

// InitInfractl creates the infragraph folder structure
func InitInfractl() error {
	root := global.RootFolder
	if global.FileExists(root) {
		return fmt.Errorf("folder \"%s\" already exists. Please delete it first", root)
	}
	subFolders := []string{global.BuildSubFolder, "recipes/std", `crafts`}
	for _, subFolder := range subFolders {
		err := os.MkdirAll(filepath.Join(root, subFolder), global.FilePerm)
		if err != nil {
			return err
		}
	}

	filesToCreate := []global.StringPair{
		{"simple_ecs.yml", recipes.SIMPLE_ECS},
		{"network.yml", recipes.NETWORK},
		{"simple_lambda.yml", recipes.SIMPLE_LAMBDA}}
	// Copy the installation files
	for _, install := range filesToCreate {
		err := global.WriteStringToFile(filepath.Join(global.RootFolder, "recipes", "std", install.S1), install.S2)
		if err != nil {
			return err
		}
	}

	fmt.Printf(`Created folder "%s" and subfolders`, root)
	return nil
}
