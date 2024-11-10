package global

import (
	"fmt"
	"os"
	"path/filepath"
)

// BuildAndDeploy builds the Terraform files, deploys them and collects outputs into the outputs map
func (lam *LambdaRecipe) BuildAndDeploy(craftPath string, outputs map[string]map[string]string) error {
	craftName := NameOnly(craftPath)
	netCraftPath, err := GetAbsoluteCraftPath(lam.NetworkCraftFile)
	if err != nil {
		return err
	}
	netCraftName := NameOnly(netCraftPath)
	net, err := NewNetworkCraft(netCraftPath)
	if err != nil {
		return err
	}
	err = net.BuildAndDeploy(netCraftPath, outputs)
	if err != nil {
		return err
	}
	netOutputs := outputs[netCraftName]
	for _, key := range []string{"vpc_id", "cluster_name", "public_subnet_id", "public_subnet_2_id", "private_subnet_id", "private_subnet_2_id"} {
		if netOutputs[key] == "" {
			return fmt.Errorf(`missing "%s" in network ecs (%s)`, key, netCraftPath)
		}
	}
	// Pick up network data
	lam.Network.ClusterName = netOutputs["cluster_name"]
	lam.Network.VpcId = netOutputs["vpc_id"]
	lam.Network.PublicSubnetId = netOutputs["public_subnet_id"]
	lam.Network.PublicSubnet2Id = netOutputs["public_subnet_2_id"]
	lam.Network.PrivateSubnetId = netOutputs["private_subnet_id"]
	lam.Network.PrivateSubnet2Id = netOutputs["private_subnet_2_id"]

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
	conf.Recipe = lam
	conf.BuildFolder = buildFolder
	conf.OutputsMap = outputs
	lam.SimpleLambda.BuildFolder = filepath.ToSlash(buildFolder)

	//lam.SetAwsEnvs()
	//if !VerifyAWSCredentials() {
	//	err = SetEnvsFromFiles(lam.InfraEnvFiles) // infra credentials
	//	if err != nil {
	//		return err
	//	}
	//	if !VerifyAWSCredentials() {
	//		return fmt.Errorf("AWS credentials not set")
	//	}
	//	lam.Creds.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	//	lam.Creds.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	//	lam.Creds.DefaultRegion = os.Getenv("AWS_DEFAULT_REGION")
	//}

	err = conf.BuildLambda()
	if err != nil {
		return err
	}
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
