package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mariusmatioc/infractl/pkg/global"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func ListClusters(craftPath string) error {
	if craftPath != "" {
		err := global.SetOsEnvsFromCraft(craftPath)
		if err != nil {
			return err
		}
	}
	return listClustersEcs()
}

type Clusters struct {
	Clusters []string `json:"clusters"`
}

func listClustersEcs() (err error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	client := ecs.NewFromConfig(cfg)
	output, err := client.ListClusters(ctx, &ecs.ListClustersInput{})
	if err != nil {
		return
	}
	result, err := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: output.ClusterArns})
	if err != nil {
		return
	}

	var clusters Clusters
	for _, cluster := range result.Clusters {
		clusters.Clusters = append(clusters.Clusters, *cluster.ClusterName)
	}
	out, _ := json.MarshalIndent(clusters, "", "    ")
	fmt.Println(string(out))
	return nil
}

//func listClustersEks() (err error) {
//	ctx := context.TODO()
//	cfg, err := config.LoadDefaultConfig(ctx)
//	if err != nil {
//		return
//	}
//
//	svc := eks.NewFromConfig(cfg)
//	output, err := svc.ListClusters(ctx, &eks.ListClustersInput{})
//	if err != nil {
//		return
//	}
//
//	var clusters Clusters
//	for _, cluster := range output.Clusters {
//		clusters.Clusters = append(clusters.Clusters, cluster)
//	}
//	out, _ := json.MarshalIndent(clusters, "", "    ")
//	fmt.Println(string(out))
//	return nil
//}
