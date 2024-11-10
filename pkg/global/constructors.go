package global

import (
	"fmt"
	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	"path/filepath"
	"strconv"
	"strings"
)

// newProject imports the docker compose file(s)
func newProject(composeFiles []string) (project *types.Project, err error) {
	projectOptions, err := cli.NewProjectOptions(composeFiles, cli.WithOsEnv,
		cli.WithInterpolation(true))
	if err != nil {
		return
	}
	err = cli.WithDotEnv(projectOptions)
	if err != nil {
		return
	}
	project, err = cli.ProjectFromOptions(projectOptions)
	if err == nil {
		err = validateProject(project)
	}
	return
}

func validateProject(project *types.Project) (err error) {
	if len(project.Services) == 0 {
		err = fmt.Errorf(`the docker_compose file does not contain any services (%s)`, project.ComposeFiles[0])
		return
	}
	for _, service := range project.Services {
		if len(service.Ports) > 1 {
			err = fmt.Errorf(`"ports": only one port mapping currently supported, service: %s`, service.Name)
			return
		}
		if len(service.Secrets) != 0 {
			err = fmt.Errorf(`"secrets" are not supported. Please use "env_file" instead, service: %s`, service.Name)
			return
		}
	}

	return
}

func NewConfig(craft *EcsRecipe) (*Config, error) {
	config := Config{}
	config.Recipe = craft
	var err error
	// Read docker compose file
	config.Compose, err = newProject(craft.DockerComposeFiles)
	return &config, err
}

/*
var dbTemplates = []DataBaseTemplate{
	{[]string{"postgres", "postgis"}, "postgres", []string{"POSTGRES_DB", "POSTGRES_NAME"}, "POSTGRES_USER", "POSTGRES_PASSWORD"},
	{[]string{"mysql", "mariadb"}, "mysql", []string{"MYSQL_DATABASE"}, "MYSQL_USER", "MYSQL_PASSWORD"},
}

func (serv *Service) DetectDatabase(compose *types.ServiceConfig) bool {
	for _, templ := range dbTemplates {
		found := false
		for _, imagePref := range templ.ImageName {
			if strings.HasPrefix(compose.Image, imagePref) {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		var s *string
		for _, name := range templ.DatabaseName {
			if s = compose.Environment[name]; s != nil {
				serv.DbName = *s
				break
			}
		}
		if serv.DbName == "" {
			continue
		}
		if s = compose.Environment[templ.Password]; s == nil {
			continue
		}
		serv.Password = *s
		if s = compose.Environment[templ.UserName]; s == nil {
			continue
		}
		serv.UserName = *s
		serv.DbEngine = templ.RdsEngineName
		return true
	}
	return false
}
*/

func FilterServices(services Services, filter func(service *Service) bool) Services {
	filtered := Services{}
	for _, serv := range services {
		if filter(serv) {
			filtered = append(filtered, serv)
		}
	}
	return filtered
}

// getHostPort decodes a string of format host:port
//func getHostPort(s string) (hostPort []string) {
//	hostPort = strings.Split(s, ":")
//	return
//}

func (config *Config) ProcessAllServices() (fargate Services, rdsMap map[string]Rdb, mqsMap map[string]Mq, err error) {
	fargate, err = config.getAllServices()
	if err != nil {
		return
	}

	config.NameToService = make(map[string]*Service)
	for _, serv := range fargate {
		config.NameToService[serv.Name] = serv
	}
	for _, serv := range fargate {
		ecs := config.Recipe.(*EcsRecipe)
		err = serv.loadComposeEnvs(filepath.Dir(ecs.DockerComposeFiles[0]))
		if err != nil {
			return
		}
		err = serv.ProcessImage()
		if err != nil {
			return
		}
	}

	// Separate Fargate from RDS and MQs
	fargate, rdsMap, mqsMap, err = config.getManagedServices(fargate)
	if err != nil {
		return
	}

	// Update the host names which in the docker compose file are the service names with the proper AWS endpoint
	for _, svc := range fargate {
		{
			svcNames := make(map[string]bool)
			for key := range config.NameToService {
				svcNames[key] = true
			}
			err = svc.updateCrossServiceHostEnvs(svcNames, "${aws_lb.%s.dns_name}")
			if err != nil {
				return
			}
		}
		if len(rdsMap) != 0 {
			svcNames := make(map[string]bool)
			for key := range rdsMap {
				svcNames[key] = true
			}
			err = svc.updateCrossServiceHostEnvs(svcNames, "${aws_db_instance.%s.address}")
			if err != nil {
				return
			}
		}
		if len(mqsMap) != 0 {
			svcNames := make(map[string]bool)
			for key := range mqsMap {
				svcNames[key] = true
			}
			err = svc.updateCrossServiceHostEnvs(svcNames, "${aws_mq_broker.%s.instances.0.endpoints.0}")
			if err != nil {
				return
			}
		}
		err = svc.ComputeDependsOn(config.NameToService, rdsMap, mqsMap)
		if err != nil {
			return
		}
	}
	return
}

// PortsFromString constructs Ports from a colon separate pair
func PortsFromString(s string) (Ports, error) {
	ports := Ports{}
	parts := strings.Split(s, ":")
	if len(parts) > 2 {
		err := fmt.Errorf("invalid port mapping: %s", s)
		return ports, err
	}
	portArr := []int{}
	for _, part := range parts {
		i, err := strconv.Atoi(part)
		if err != nil {
			err = fmt.Errorf("invalid port mapping: %s", s)
			return ports, err
		}
		portArr = append(portArr, i)
	}
	ports.Target = portArr[len(portArr)-1]
	ports.Published = portArr[0]
	return ports, nil
}
