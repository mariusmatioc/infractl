package global

import (
	"github.com/compose-spec/compose-go/types"
)

type CraftSection struct {
	RecipeVersion      string        `yaml:"recipe_version"`
	Type               string        `yaml:"type"` // ecs, network, lambda
	InfraEnvFile       string        `yaml:"infra_env_file"`
	DockerComposeFiles []string      `yaml:"docker_compose_files"`
	NetworkCraftFile   string        `yaml:"network_craft_file"`
	Name               string        // The name of the craft file
	BackendConfig      BackendConfig // optional
}

type BackendConfig struct {
	Bucket string
	Key    string
}

type Network struct {
	ClusterName      string `yaml:"cluster_name"`
	VpcCidr          string `yaml:"vpc_cidr"`
	VpcId            string `yaml:"vpc_id"`
	PublicSubnetId   string
	PublicSubnet2Id  string
	PrivateSubnetId  string
	PrivateSubnet2Id string
}

type SimpleEcs struct {
	Type        string                        `yaml:"type"`         // fargate or EC2
	MachineType string                        `yaml:"machine_type"` // only for EC2
	Services    map[string]RecipeServiceItems `yaml:"services"`
}

type SimpleRds struct {
	Databases map[string]DbItems `yaml:"services"`
}

// SimpleMqs is for AWS managed message queues
type SimpleMqs struct {
	Queues map[string]MqItems `yaml:"services"`
}

type SimpleLambda struct {
	FunctionName     string `yaml:"function_name"`
	Handler          string `yaml:"handler"`
	Runtime          string `yaml:"runtime"`
	SourceFolder     string `yaml:"source_folder"`
	MemorySize       int    `yaml:"memory_size"`
	Timeout          int    `yaml:"timeout"`
	EphemeralStorage int    `yaml:"ephemeral_storage"`
	LambdaTriggers   `yaml:"triggers"`
	Environment      map[string]string `yaml:"environment"`
	EnvFiles         []string          `yaml:"env_files"`
	Layers           []string          `yaml:"layers"`

	EnvsString   string // Constructed
	LayersString string // Constructed
	BuildFolder  string // Constructed
}

type LambdaTriggers struct {
	ScheduleExpression string `yaml:"schedule_expression"`
	S3ObjectCreated    string `yaml:"s3_object_created"`
}

type External struct {
	Name      string `yaml:"name"`
	CraftFile string `yaml:"craft_file"` // If relative, then to the location of enclosing file
}

type DbItems struct {
	Public      bool     `yaml:"public"`
	DbEngine    string   `yaml:"engine"`
	MachineType string   `yaml:"machine_type"`
	StorageType string   `yaml:"storage_type"`
	StorageIops int      `yaml:"storage_iops"`
	StorageGigs int      `yaml:"storage_gigs"`
	EnvFiles    []string `yaml:"env_files"`
}

type MqItems struct {
	Public   bool     `yaml:"public"`
	Engine   string   `yaml:"engine"`
	EnvFiles []string `yaml:"env_files"`
}

// RecipeServiceItems is the configuration for each service from the recipe
type RecipeServiceItems struct {
	DesiredNodes      int               `yaml:"desired_nodes"`
	Cpu               int               `yaml:"cpu"`
	Memory            int               `yaml:"memory"`
	DomainName        string            `yaml:"domain_name"`
	LoadBalancerHttp  []string          `yaml:"load_balancer_http"`
	LoadBalancerHttps []string          `yaml:"load_balancer_https"`
	Environment       map[string]string `yaml:"environment"`
	EnvFiles          []string          `yaml:"env_files"`
	DependsOn         []string          `yaml:"depends_on"`
}

type Creds struct {
	AccessKey     string `yaml:"AWS_ACCESS_KEY_ID"`
	SecretKey     string `yaml:"AWS_SECRET_ACCESS_KEY"`
	DefaultRegion string `yaml:"AWS_DEFAULT_REGION"`
}

// Config is top level collection of configuration info
type Config struct {
	Compose *types.Project // From compose file
	Recipe  any            // ECSRecipe or NetworkRecipe

	NameToService map[string]*Service
	BuildFolder   string
	// The Terraform outputs to be used as environment variables in other crafts
	// The first key is the craft name, the second key is the output name
	OutputsMap map[string]map[string]string
}

func (cfg *Config) GetEcsRecipe() *EcsRecipe {
	return cfg.Recipe.(*EcsRecipe)
}

func (cfg *Config) GetNetworkRecipe() *NetworkRecipe {
	return cfg.Recipe.(*NetworkRecipe)
}

func (cfg *Config) GetLambdaRecipe() *LambdaRecipe {
	return cfg.Recipe.(*LambdaRecipe)
}

// Ports is mapping the ports on the task or load balancer
// For Fargate task, both bort must be the same
type Ports struct {
	Target    int // On the container
	Published int // On the host
}

type HealthCheck struct {
	Test     string
	Interval int
	Timeout  int
	Retries  uint
}

type Service struct {
	Name                 string
	Ports                []Ports
	DbPort               int
	LoadBalancerPortsTCP []Ports
	LoadBalancerPortsTLS []Ports
	LoadBalacerTargets   []int
	NeedsCertificate     bool
	Memory               int // in MiB
	Cpu                  int // 1024 = 1vCPU
	DesiredCount         int
	EnvsString           string
	Entrypoint           string
	IsOnlyService        bool // true if there is only one service in the compose file

	BuildContext    string // folder of where the build is done
	BuildHash       string // hash of the build context
	TaskHash        string // the last 10 characters of BuildHash
	BuildDockerfile string // name of Dockerfile
	BuildTarget     string

	Image     string // either a docker registry image or an ECR image after build and push
	DependsOn string

	HealthCheck HealthCheck
	DomainName  string
	HostedZone  string

	composeConfig types.ServiceConfig // From compose file
	envMap        map[string]string
}

// Rdb is a managed relational database
type Rdb struct {
	Name        string // From service name in compose file
	DbEngine    string
	DbName      string
	UserName    string
	Password    string
	Port        int
	Public      bool
	MachineType string
	StorageType string
	StorageIops int
	StorageGigs int
}

// Mq is a managed message queue
type Mq struct {
	Name     string // From service name in compose file
	Engine   string // from craft file
	Public   bool   // from craft file
	UserName string // from envs
	Password string // from envs
}

type Services []*Service

type DataBaseTemplate struct {
	ImageName     []string
	RdsEngineName string
	DatabaseName  []string
	UserName      string
	Password      string
}

type StringPair struct {
	S1 string
	S2 string
}
