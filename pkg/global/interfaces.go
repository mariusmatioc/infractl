package global

type Craft interface {
	GetCraftSection() *CraftSection
	GetCreds() *Creds
}

func (ecs *EcsRecipe) GetCraftSection() *CraftSection {
	return &ecs.CraftSection
}

func (net *NetworkRecipe) GetCraftSection() *CraftSection {
	return &net.CraftSection
}

func (lam *LambdaRecipe) GetCraftSection() *CraftSection {
	return &lam.CraftSection
}

func (ecs *EcsRecipe) GetCreds() *Creds {
	return &ecs.Creds
}

func (net *NetworkRecipe) GetCreds() *Creds {
	return &net.Creds
}

func (lam *LambdaRecipe) GetCreds() *Creds {
	return &lam.Creds
}
