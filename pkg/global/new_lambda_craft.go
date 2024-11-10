package global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func newLambdaCraft(data []byte, path string) (*LambdaRecipe, error) {
	var lam LambdaRecipe
	err := yaml.Unmarshal(data, &lam)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %s", path, err.Error())
	}
	sl := &lam.SimpleLambda
	mandatoryStringFields := []string{sl.FunctionName, sl.Handler, sl.Runtime, sl.SourceFolder}
	mandatoryStringFieldNames := []string{"function_name", "handler", "runtime", "source_folder"}
	for i, field := range mandatoryStringFields {
		if field == "" {
			return nil, fmt.Errorf(`missing "%s" in (%s)`, mandatoryStringFieldNames[i], path)
		}
	}
	sl.SourceFolder = filepath.ToSlash(ToAbsPathBasedOn(RootParent, sl.SourceFolder))
	if sl.EphemeralStorage == 0 {
		sl.EphemeralStorage = 512 // default
	}
	// Triggers
	if sl.ScheduleExpression == "" && sl.S3ObjectCreated == "" {
		return nil, fmt.Errorf(`missing event trigger in (%s)`, path)
	}
	if len(sl.Layers) > 0 {
		sl.LayersString = fmt.Sprintf(`"%s"`, strings.Join(sl.Layers, `",\n"`))
	}

	// Environment files
	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	err = os.Chdir(RootParent) // we want relative file names to refer to parent of root folder
	if err != nil {
		return nil, err
	}

	// Env vars
	err = sl.createEnvsString()
	if err != nil {
		return nil, err
	}
	return &lam, nil
}

func (sl *SimpleLambda) createEnvsString() (err error) {
	envMap := make(map[string]string)
	for _, filePath := range sl.EnvFiles {
		filePath = os.ExpandEnv(filePath)
		err = ReadEnvFile(filePath, envMap)
		if err != nil {
			return
		}
	}
	for key, val := range sl.Environment {
		name := strings.TrimSpace(key)
		value := strings.TrimSpace(val)
		envMap[name] = value
	}
	envs := []string{}
	for key, val := range envMap {
		envs = append(envs, fmt.Sprintf(`     %s = "%s"`, key, val))
	}
	sl.EnvsString = strings.Join(envs, "\n")
	return nil
}
