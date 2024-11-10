package global

import (
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// NewCraft returns a ECSRecipe, NetworkRecipe or LambdaRecipe
func NewCraft(path string) (Craft, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Create a struct to hold the YAML data
	type CraftOnly struct {
		CraftSection `yaml:"craft"`
	}
	var craftSection CraftOnly
	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &craftSection)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %s", path, err.Error())
	}

	var craft Craft
	switch craftSection.Type {
	case "ecs":
		craft, err = newECSCraft(data, path)
	case "network":
		craft, err = newNetworkCraft(data, path)
	case "lambda":
		craft, err = newLambdaCraft(data, path)
	default:
		err = fmt.Errorf(`unknown "craft type": "%s" in (%s)`, craftSection.Type, path)
	}
	if err != nil {
		return nil, err
	}
	cs := craft.GetCraftSection()
	craftName := NameOnly(path)
	cs.Name = craftName
	cs.ExpandEnvsInFileNames()
	infraEnvs, err := GetEnvMapFromFiles([]string{cs.InfraEnvFile})
	if Backend != nil {
		// Using remote backend
		cs.BackendConfig = *Backend
	}

	creds := craft.GetCreds()
	err = adjustCredentialItem(&creds.DefaultRegion, "AWS_DEFAULT_REGION", infraEnvs)
	if err != nil {
		return nil, err
	}
	err = adjustCredentialItem(&creds.AccessKey, "AWS_ACCESS_KEY_ID", infraEnvs)
	if err != nil {
		return nil, err
	}
	err = adjustCredentialItem(&creds.SecretKey, "AWS_SECRET_ACCESS_KEY", infraEnvs)
	return craft, err
}

func newNetworkCraft(data []byte, path string) (*NetworkRecipe, error) {
	var net NetworkRecipe
	err := yaml.Unmarshal(data, &net)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %s", path, err.Error())
	}
	if net.Network.VpcCidr == "" && net.Network.VpcId == "" {
		return nil, fmt.Errorf(`you must specify one of "vpc_id" or "vpc_cidr" in (%s)`, path)
	}
	if net.Network.ClusterName == "" {
		return nil, fmt.Errorf(`missing "cluster_name" in (%s)`, path)
	}
	if net.NetworkCraftFile != "" {
		return nil, fmt.Errorf(`"network_craft_file" not allowed in a craft of type network (%s)`, path)
	}
	// Environment files
	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	err = os.Chdir(RootParent) // we want relative file names to refer to parent of root folder
	if err != nil {
		return nil, err
	}
	net.Name = NameOnly(path)

	return &net, nil
}

func NewNetworkCraft(path string) (*NetworkRecipe, error) {
	craft, err := NewCraft(path)
	if err != nil {
		return nil, err
	}
	if net, ok := craft.(*NetworkRecipe); ok {
		return net, nil
	}
	return nil, fmt.Errorf(`craft file %s is not of type network`, path)
}

func newECSCraft(data []byte, craftPath string) (*EcsRecipe, error) {
	var ecs EcsRecipe
	err := yaml.Unmarshal(data, &ecs)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %s", craftPath, err.Error())
	}
	ecs.SimpleEcs.Type = AdjustAwsString(ecs.SimpleEcs.Type)
	if ecs.NetworkCraftFile == "" {
		return nil, fmt.Errorf(`"missing "network_craft_file" in (%s)`, craftPath)
	}
	if ecs.SimpleEcs.Type != "fargate" {
		return nil, fmt.Errorf(`only "fargate" is currently supported (%s)`, craftPath)
	}
	if len(ecs.DockerComposeFiles) == 0 {
		return nil, fmt.Errorf(`missing "docker_compose_files" in (%s)`, craftPath)
	}
	return &ecs, nil
}

//func (ecs *EcsRecipe) GetSimpleEcss(craftPath string) ([]SimpleEcs, error) {
//	ecss := []SimpleEcs{ecs.SimpleEcs}
//	for _, external := range ecs.Externals {
//		extPath := ToAbsPathBasedOn(filepath.Dir(craftPath), external.CraftFile)
//		ext, err := NewEcsCraft(extPath)
//		if err != nil {
//			return nil, err
//		}
//		ecss = append(ecss, ext.SimpleEcs)
//	}
//	return ecss, nil
//}

func NewEcsCraft(path string) (*EcsRecipe, error) {
	craft, err := NewCraft(path)
	if err != nil {
		return nil, err
	}
	if ecs, ok := craft.(*EcsRecipe); ok {
		return ecs, nil
	}
	return nil, fmt.Errorf(`craft file %s is not of type ecs`, path)
}

func getLoadBalancerPorts(lbPortStrings []string, svcPorts []Ports, svcName string) (lbPorts []Ports, err error) {
	for _, lbPortStr := range lbPortStrings {
		ports := Ports{}
		ports, err = PortsFromString(lbPortStr)
		if err != nil {
			err = fmt.Errorf("load balancer port syntax error for service %s (%s): %s", svcName, lbPortStr, err.Error())
			return
		}
		lbPorts = append(lbPorts, ports)
		found := false
		for _, svpPort := range svcPorts {
			if ports.Target == svpPort.Published {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf(`load balancer port %d for service "%s" has no match in docker compose file`, ports.Target, svcName)
			return
		}
	}
	return
}

func SetOsEnvsFromCraft(craftPath string) (err error) {
	craft, err := NewCraft(craftPath)
	if err != nil {
		return
	}
	craft.GetCreds().SetOsEnvs()
	return
}

func (creds *Creds) SetOsEnvs() {
	_ = os.Setenv("AWS_DEFAULT_REGION", creds.DefaultRegion)
	_ = os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKey)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretKey)
}

func ExpandEnvInFileNames(files []string) {
	for i := range files {
		files[i] = os.ExpandEnv(files[i])
	}
}

func (cs *CraftSection) ExpandEnvsInFileNames() {
	ExpandEnvInFileNames(cs.DockerComposeFiles)
	cs.InfraEnvFile = os.ExpandEnv(cs.InfraEnvFile)
}

func CollectCraftNames(path string, names map[string]bool) (err error) {
	craft, err := NewCraft(path)
	if err != nil {
		return
	}
	names[filepath.Base(path)] = true
	if ecs, ok := craft.(*EcsRecipe); ok {
		names[craft.GetCraftSection().NetworkCraftFile] = true
		for _, external := range ecs.Externals {
			extPath := ToAbsPathBasedOn(filepath.Dir(path), external.CraftFile)
			err = CollectCraftNames(extPath, names)
			if err != nil {
				return
			}
		}
	}
	return
}

// adjustCredentialItem adjust existing item from envs
func adjustCredentialItem(item *string, itemName string, envs map[string]string) error {
	if *item == "" {
		if val, ok := envs[itemName]; ok {
			*item = val
		} else {
			return fmt.Errorf(`missing credential "%s"`, itemName)
		}
	}
	// Otherwise, it is a literal value
	return nil
}
