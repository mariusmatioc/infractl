package global

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func ListServices(craftPath string) error {
	ecsCraft, err := NewEcsCraft(craftPath)
	if err != nil {
		return err
	}
	netCraftName := NameOnly(ecsCraft.NetworkCraftFile)
	netBuildFolder, err := GetBuildFolder(netCraftName)
	outs := make(map[string]string)
	err = GetTerraformOutputs(netBuildFolder, outs)
	if err != nil {
		return err
	}
	ecsCraft.Network.ClusterName = outs["cluster_name"]
	if ecsCraft.Network.ClusterName == "" {
		return fmt.Errorf("missing cluster_name in network craft (%s)", netCraftName)
	}

	ecsCraft.GetCreds().SetOsEnvs()
	clusterMap := make(map[string]bool)
	cluster := ecsCraft.Network.ClusterName
	clusterMap[cluster] = true
	err = listServicesEcs(cluster)
	if err != nil {
		return err
	}
	return nil
}

type EcsService struct {
	Name         string `json:"name"`
	Arn          string `json:"arn"`
	Status       string `json:"status"`
	DesiredCount int    `json:"desired_count"`
	RunningCount int    `json:"running_count"`
}

func listServicesEcs(cluster string) (err error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	svc := ecs.NewFromConfig(cfg)
	output, err := svc.ListServices(ctx, &ecs.ListServicesInput{Cluster: &cluster})
	if err != nil {
		return
	}

	services, err := svc.DescribeServices(ctx, &ecs.DescribeServicesInput{Cluster: &cluster, Services: output.ServiceArns})
	svcs := make([]EcsService, 0)
	for _, service := range services.Services {
		svc := EcsService{
			Name:         *service.ServiceName,
			Arn:          *service.ServiceArn,
			Status:       *service.Status,
			DesiredCount: int(service.DesiredCount),
			RunningCount: int(service.RunningCount),
		}
		svcs = append(svcs, svc)
	}
	type Result struct {
		Cluster  string       `json:"cluster"`
		Services []EcsService `json:"services"`
	}
	result := Result{Cluster: cluster, Services: svcs}
	j, _ := json.MarshalIndent(result, "", "    ")
	fmt.Println(string(j))
	return nil
}

//type EksService struct {
//	Namespace  string `json:"namespace"`
//	Name       string `json:"name"`
//	Type       string `json:"type"`
//	ClusterIP  string `json:"cluster-ip"`
//	ExternalIP string `json:"external-ip,omitempty"`
//	Ports      string `json:"ports"`
//	Age        string `json:"age"`
//}

//func listServicesEks(cluster, region string) (err error) {
//	// Point kubectl to the cluster
//	_, err = exec.Command("aws", "eks", "update-kubeconfig", "--region", region, "--name", cluster).Output()
//	if err != nil {
//		return
//	}
//	out, err := utils.RunKubectlCommand([]string{"get", "svc", "-A"})
//	if err != nil {
//		return
//	}
//	scanner := bufio.NewScanner(strings.NewReader(out))
//	// Skip the first line (header)
//	scanner.Scan()
//	svcs := make([]EksService, 0)
//	for scanner.Scan() {
//		parts := strings.Fields(strings.TrimSpace(scanner.Text()))
//		service := EksService{}
//		service.Namespace = parts[0]
//		service.Name = parts[1]
//		service.Type = parts[2]
//		service.ClusterIP = parts[3]
//		if parts[4] != "<none>" {
//			service.ExternalIP = parts[4]
//		}
//		service.Ports = parts[5]
//		service.Age = parts[6]
//
//		svcs = append(svcs, service)
//	}
//	if err = scanner.Err(); err != nil {
//		return
//	}
//	j, _ := json.MarshalIndent(svcs, "", "    ")
//	fmt.Println(string(j))
//	return
//}
