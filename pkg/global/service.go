package global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/types"
	"github.com/google/uuid"
)

// Creates a Service from the docker compose info
func newService(compose types.ServiceConfig) (*Service, error) {
	ecs := Service{}
	ecs.envMap = make(map[string]string)
	ecs.Name = strings.ReplaceAll(compose.Name, "_", "-")
	ecs.composeConfig = compose

	// ------------------- Ports ----------------------------------------------------------------------------
	for _, p := range compose.Ports {
		//pub, _ := strconv.Atoi(p.Published)
		//if pub == 0 {
		//	pub = int(p.Target)
		//}
		//pub = int(p.Target)
		// For ECS, the ports must be the same
		ports := Ports{
			Target:    int(p.Target),
			Published: int(p.Target),
		}
		ecs.Ports = append(ecs.Ports, ports)
	}
	if len(ecs.Ports) != 0 {
		ecs.DbPort = ecs.Ports[0].Target
	}

	// ----------------------------------------- Volumes -------------------------------------------------
	if len(compose.Volumes) > 0 {
		fmt.Println("Volumes are not supported in Fargate. Please remove. Add copy statements to the docker file if needed.", ecs.Name)
	}

	// --------------------------------------------- Command ----------------------------------------------
	commands := []string{}
	if len(compose.Command) != 0 {
		for _, e := range compose.Command {
			commands = append(commands, `"`+e+`"`)
		}
	} else if len(compose.Entrypoint) != 0 {
		for _, e := range compose.Entrypoint {
			commands = append(commands, `"`+e+`"`)
		}
	}
	if len(commands) != 0 {
		str := strings.Join(commands, ", ")
		ecs.Entrypoint = fmt.Sprintf("[%s]", str)
	}

	// --------------------------------------------- HealthCheck ----------------------------------------------
	if compose.HealthCheck != nil && len(compose.HealthCheck.Test) != 0 {
		ecs.HealthCheck.Test = QuotedArray(compose.HealthCheck.Test)
		if compose.HealthCheck.Interval != nil {
			ecs.HealthCheck.Interval = int(time.Duration(*compose.HealthCheck.Interval).Seconds())
		}
		if compose.HealthCheck.Timeout != nil {
			ecs.HealthCheck.Interval = int(time.Duration(*compose.HealthCheck.Timeout).Seconds())
		}
		if compose.HealthCheck.Retries != nil {
			ecs.HealthCheck.Retries = uint(*compose.HealthCheck.Retries)
		}
	}

	// NOTE: Some fields will be processed/updated later
	return &ecs, nil
}

// getAllServices returns all the services from the compose files, separating database services
func (config *Config) getAllServices() (services Services, err error) {
	publishedPorts := make(map[int]bool)
	for _, serv := range config.Compose.Services {
		sv, err2 := newService(serv)
		if err2 != nil {
			err = err2
			return
		}
		for _, ports := range sv.Ports {
			if _, ok := publishedPorts[ports.Published]; ok {
				err = fmt.Errorf("duplicate published port found: %d", ports.Published)
				return
			}
			publishedPorts[ports.Published] = true
		}
		services = append(services, sv)
	}
	// Make sure that all services in recipe exist in compose
	for name := range config.GetEcsRecipe().SimpleEcs.Services {
	        name_formatted := strings.ReplaceAll(name, "_", "-")  // Must use the same name processing as in newService() function to match the sv.Name in "services"
		found := false
		for _, sv := range services {
			if sv.Name == name_formatted {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("service %s from recipe is not defined in docker-compose file", name)
			return
		}
	}
	// Update from recipe
	for _, sv := range services {
		err = sv.UpdateFromRecipe(config)
		if err != nil {
			break
		}
	}
	return
}

func (serv *Service) ProcessImage() error {
	composeBuild := serv.composeConfig.Build
	if composeBuild != nil {
		serv.Image = fmt.Sprintf("aws_ecr_repository.%s.repository_url", serv.Name)
		serv.BuildContext = filepath.ToSlash(composeBuild.Context)
		if ForceRebuild {
			// We will use a uuid instead of a hash to force a new build
			id := uuid.NewString()
			serv.BuildHash = strings.ReplaceAll(id, "-", "")
		} else {
			buildHash, err := HashOfFolder(serv.BuildContext)
			if err != nil {
				return err
			}
			serv.BuildHash = buildHash
		}
		serv.TaskHash = serv.BuildHash[:10]
		serv.BuildDockerfile = composeBuild.Dockerfile
		serv.BuildTarget = composeBuild.Target
	} else {
		serv.Image = fmt.Sprintf(`"%s"`, serv.composeConfig.Image)
		serv.TaskHash = "task"
	}
	return nil
}

func (serv *Service) ComputeDependsOn(nameToService map[string]*Service, rdsMap map[string]Rdb, mqsMap map[string]Mq) error {
	var dependsOn []string
	composeBuild := serv.composeConfig.Build
	if composeBuild != nil {
		dependsOn = append(dependsOn, fmt.Sprintf("docker_image.%s_image", serv.Name))
	} else {
	}

	for dep := range serv.composeConfig.DependsOn {
		if _, ok := nameToService[dep]; ok {
			dependsOn = append(dependsOn, "aws_ecs_task_definition."+dep)
		} else if _, ok := rdsMap[dep]; ok {
			dependsOn = append(dependsOn, "aws_db_instance."+dep)
		} else if _, ok := mqsMap[dep]; ok {
			dependsOn = append(dependsOn, "aws_mq_broker."+dep)
		} else {
			return fmt.Errorf("service %s depends on %s which is not defined in docker-compose file", serv.Name, dep)
		}
	}
	if len(dependsOn) > 0 {
		serv.DependsOn = fmt.Sprintf("depends_on = [%s]", strings.Join(dependsOn, ", "))
	}
	return nil
}

func (srv *Service) CreateEnvsString() {
	envs := []string{}
	for key, val := range srv.envMap {
		envs = append(envs, envToString(key, val))
	}
	srv.EnvsString = fmt.Sprintf("[%s]", strings.Join(envs, ","))
}

// UpdateFromRecipe updates the service from new data in the recipe
func (svc *Service) UpdateFromRecipe(cfg *Config) error {
	// find service, if any
	ecs := cfg.Recipe.(*EcsRecipe)
	if prov, ok := ecs.SimpleEcs.Services[svc.Name]; !ok {
		// defaults
		svc.DesiredCount = 1
		svc.Memory, svc.Cpu, _ = bestMemCpuMatch(0, 0)
		svc.LoadBalancerPortsTCP = svc.Ports
	} else {
		wd, _ := os.Getwd()
		defer func() { _ = os.Chdir(wd) }()
		err := os.Chdir(RootParent) // we want relative file names to refer to parent of root folder
		if err != nil {
			return err
		}
		// --------------------------------------------- Load Balancer ----------------------------------------------
		// Check that ports in recipe and docker compose match
		lbPortsTCP := []Ports{}
		lbPortsTCP, err = getLoadBalancerPorts(prov.LoadBalancerHttp, svc.Ports, svc.Name)
		if err != nil {
			return err
		}
		// all good
		svc.LoadBalancerPortsTCP = lbPortsTCP
		lbPortsTLS := []Ports{}
		lbPortsTLS, err = getLoadBalancerPorts(prov.LoadBalancerHttps, svc.Ports, svc.Name)
		if err != nil {
			return err
		}
		// all good
		svc.LoadBalancerPortsTLS = lbPortsTLS
		svc.NeedsCertificate = len(svc.LoadBalancerPortsTLS) > 0
		// --------------------------------------------- DomainName ----------------------------------------------
		svc.DomainName = prov.DomainName
		if svc.DomainName != "" {
			// Get top hosted zone
			parts := strings.Split(svc.DomainName, ".")
			svc.HostedZone = strings.Join(parts[len(parts)-2:], ".")
		}
		// --------------------------------------------- Desired count ----------------------------------------------
		svc.DesiredCount = prov.DesiredNodes
		if svc.DesiredCount <= 0 {
			svc.DesiredCount = 1
		}
		// ----------------------------------------- Memory and CPU -------------------------------------------------
		svc.Memory, svc.Cpu, err = bestMemCpuMatch(prov.Memory, prov.Cpu)
		if err != nil {
			return err
		}
		// ----------------------------------------- Environment -------------------------------------------------
		for _, filePath := range prov.EnvFiles {
			filePath = os.ExpandEnv(filePath)
			err = svc.ReadEnvFile(filePath)
			if err != nil {
				return err
			}
		}
		// Map the specified env values, while replacing externals
		for key, val := range prov.Environment {
			name := strings.TrimSpace(key)
			value := strings.TrimSpace(val)
			if strings.HasPrefix(value, "external.") {
				return fmt.Errorf(`foound "external." in env variable %s in service %s. Use "externals." instead`, name, svc.Name)
			}
			if strings.HasPrefix(value, "externals.") {
				parts := strings.Split(value, ".")
				if len(parts) != 3 {
					err := fmt.Errorf("invalid externals reference %s in service %s", value, svc.Name)
					return err
				}
				externalName := parts[1]
				externalKey := parts[2]
				outputMap, ok := cfg.OutputsMap[externalName]
				if !ok {
					err := fmt.Errorf("invalid external reference %s in env variable %s service %s", externalName, name, svc.Name)
					return err
				}
				externalValue, ok := outputMap[externalKey]
				if !ok {
					err := fmt.Errorf("there is no output called %s in external %s", externalKey, externalName)
					return err
				}
				value = externalValue
			}
			svc.envMap[name] = value
		}
	}
	// Post process load balancer ports
	published := make(map[int]bool)
	target := make(map[int]bool)
	for _, ports := range svc.LoadBalancerPortsTCP {
		if published[ports.Published] {
			err := fmt.Errorf("duplicate published port %d for service %s", ports.Published, svc.Name)
			return err
		}
		published[ports.Published] = true
		target[ports.Target] = true
	}
	for _, ports := range svc.LoadBalancerPortsTLS {
		if published[ports.Published] {
			err := fmt.Errorf("duplicate published port %d for service %s", ports.Published, svc.Name)
			return err
		}
		published[ports.Published] = true
		target[ports.Target] = true
	}
	svc.LoadBalacerTargets = []int{}
	for k := range target {
		svc.LoadBalacerTargets = append(svc.LoadBalacerTargets, k)
	}
	return nil
}

//func (serv *Service) IsDb() bool {
//	return serv.DbEngine != ""
//}
