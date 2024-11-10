package global

import (
	"fmt"
	"os"
)

// BuildAndDeploy builds the Terraform files, deploys them and collects outputs into the outputs map
func (ecs *EcsRecipe) BuildAndDeploy(craftPath, externalName string, outputs map[string]map[string]string) error {
	// Build the network ecs first, as the main ecs depends on it
	netCraftPath, err := GetAbsoluteCraftPath(ecs.NetworkCraftFile)
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

	// Now this ecs
	craftName := NameOnly(craftPath)
	// Pick up network data
	ecs.Network.ClusterName = netOutputs["cluster_name"]
	ecs.Network.VpcId = netOutputs["vpc_id"]
	ecs.Network.PublicSubnetId = netOutputs["public_subnet_id"]
	ecs.Network.PublicSubnet2Id = netOutputs["public_subnet_2_id"]
	ecs.Network.PrivateSubnetId = netOutputs["private_subnet_id"]
	ecs.Network.PrivateSubnet2Id = netOutputs["private_subnet_2_id"]

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

	conf, err := NewConfig(ecs)
	if err != nil {
		return err
	}
	conf.BuildFolder = buildFolder
	conf.OutputsMap = outputs

	err = conf.BuildEcs()
	if err != nil {
		return err
	}

	err = TerraformDeploy(buildFolder, craftName)
	if err != nil {
		return err
	}

	// Collect outputs
	if externalName != "" {
		outs := make(map[string]string)
		err = GetTerraformOutputs(buildFolder, outs)
		if err != nil {
			return err
		}
		fmt.Printf("'%s' had %d outputs\n", externalName, len(outs))
		if len(outs) > 0 {
			outputs[externalName] = outs
		}
	}
	return nil
}

// {{if .IsOnlyService}}
//	}
//{{else}}
//      dependsOn: [{condition: "SUCCESS", containerName: "{{.Name}}_Sidecar"}],
//    },
//    {
//     name : "{{.Name}}_Sidecar",
//     image : "docker/ecs-searchdomain-sidecar:1.0",
//     command: [format("%s.compute.internal", local.region), local.cloudmap_name],
//     essential: false,
//     portMappings : [{
//       containerPort : 85,
//       protocol : "tcp"
//     }],
//    }
//{{end}}

//
///* Service discovery */
//resource "aws_service_discovery_private_dns_namespace" "cloudmap" {
//name        = local.cloudmap_name
//vpc         = data.aws_vpc.vpc.id
//}

//{{if .IsOnlyService}}
//{{else}}
//service_registries {
//registry_arn = aws_service_discovery_service.{{.Name}}.arn
//}
//{{end}}
