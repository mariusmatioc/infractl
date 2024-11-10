package global

type EcsRecipe struct {
	CraftSection `yaml:"craft"`
	SimpleEcs    SimpleEcs  `yaml:"simple_ecs"`
	SimpleRds    SimpleRds  `yaml:"simple_rds"`
	SimpleMqs    SimpleMqs  `yaml:"simple_mqs"`
	Externals    []External `yaml:"externals"`
	Creds        `yaml:"credentials"`
	Network      Network // Not in craft file
}

type NetworkRecipe struct {
	CraftSection `yaml:"craft"`
	Network      Network `yaml:"network"`
	Creds        `yaml:"credentials"`
}

type LambdaRecipe struct {
	CraftSection `yaml:"craft"`
	Network      Network      `yaml:"network"`
	SimpleLambda SimpleLambda `yaml:"simple_lambda"`
	Creds        `yaml:"credentials"`
}
