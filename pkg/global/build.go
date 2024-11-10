package global

import (
	"fmt"
	"github.com/mariusmatioc/infractl/pkg/fargate"
	"github.com/mariusmatioc/infractl/pkg/lambda"
)

// BuildEcs creates the Terraform files from the Fargate templates
// configPath is the Infragraph configuration
func (cfg *Config) BuildEcs() (err error) {
	ecs := cfg.GetEcsRecipe()
	err = cfg.BuildFromTemplate(fargate.Providers, "providers.tf", ecs)
	if err != nil {
		return
	}
	services, dbs, mqs, err := cfg.ProcessAllServices()
	if err != nil {
		return
	}
	err = cfg.writeFargateFile("data.tf", fargate.Data)
	if err != nil {
		return
	}
	err = cfg.writeFargateFile("logs.tf", fargate.Logs)
	if err != nil {
		return
	}
	err = cfg.BuildFromTemplate(fargate.Locals, "locals.tf", cfg)
	if err != nil {
		return
	}
	err = cfg.BuildFromTemplate(fargate.NetworkUse, "network.tf", ecs.Network)
	if err != nil {
		return
	}

	if len(dbs) > 0 {
		dbsSlice := []any{}
		for _, db := range dbs {
			dbsSlice = append(dbsSlice, db)
		}
		err = cfg.BuildSliceFromTemplate(fargate.Rd, "rds.tf", dbsSlice)
		if err != nil {
			return
		}
	}

	if len(mqs) > 0 {
		mqsSlice := []any{}
		for _, mq := range mqs {
			mqsSlice = append(mqsSlice, mq)
		}
		err = cfg.BuildSliceFromTemplate(fargate.Mq, "mqs.tf", mqsSlice)
		if err != nil {
			return
		}
	}

	if len(services) > 0 {
		err = cfg.writeFargateFile("iam.tf", fargate.Iam)
		if err != nil {
			return
		}

		err = cfg.BuildFromServiceTemplate(fargate.SecurityGroup, "security-groups.tf", services)
		if err != nil {
			return
		}
		err = buildFargateServices(cfg, services)
		if err != nil {
			return
		}
	}
	return
}

func buildFargateServices(cfg *Config, services Services) (err error) {
	buildServices := FilterServices(services, func(serv *Service) bool { return serv.BuildContext != "" })
	if len(buildServices) != 0 {
		err = cfg.BuildFromServiceTemplate(fargate.Ecr, "ecr.tf", buildServices)
		if err != nil {
			return
		}
	}
	for _, serv := range services {
		serv.CreateEnvsString()
	}
	if len(services) > 1 {
		//err = cfg.BuildFromServiceTemplate(fargate.ServiceDiscovery, "service-discovery.tf", services)
		//if err != nil {
		//	return
		//}
	} else {
		services[0].IsOnlyService = true
	}

	err = cfg.BuildFromServiceTemplate(fargate.Service, "services.tf", services)
	if err != nil {
		return
	}

	err = cfg.BuildFromServiceTemplate(fargate.TaskDef, "task-defs.tf", services)
	if err != nil {
		return
	}
	return
}

func (cfg *Config) BuildNetwork() (err error) {
	net := cfg.GetNetworkRecipe()
	err = cfg.BuildFromTemplate(fargate.Providers, "providers.tf", net)
	if err != nil {
		return
	}
	err = cfg.writeFargateFile("data.tf", fargate.Data)
	if err != nil {
		return
	}
	err = cfg.BuildFromTemplate(fargate.LocalsNetwork, "locals.tf", cfg)
	if err != nil {
		return
	}
	if net.Network.VpcId == "" {
		if net.Network.VpcCidr == "" {
			err = fmt.Errorf("vpc_cidr is required if vpc_id is not provided")
		} else {
			// Create the VPC
			err = cfg.BuildFromTemplate(fargate.NetworkCreate, "network.tf", net.Network)
		}
	} else {
		// Use the existing VPC
		err = cfg.BuildFromTemplate(fargate.NetworkExisting, "network.tf", net.Network)
	}
	if err != nil {
		return
	}
	return
}

func (cfg *Config) BuildLambda() (err error) {
	lam := cfg.GetLambdaRecipe()
	err = cfg.BuildFromTemplate(lambda.Providers, "providers.tf", lam)
	if err != nil {
		return
	}
	err = cfg.writeFargateFile("data.tf", lambda.Data)
	if err != nil {
		return
	}
	err = cfg.BuildFromTemplate(lambda.Locals, "locals.tf", lam)
	if err != nil {
		return
	}
	err = cfg.BuildFromTemplate(lambda.NetworkUse, "network.tf", lam.Network)
	if err != nil {
		return
	}
	err = cfg.BuildFromTemplate(lambda.Lambda, "lambda.tf", lam.SimpleLambda)
	return
}
