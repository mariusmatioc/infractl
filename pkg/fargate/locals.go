package fargate

const Locals = `
locals {
  name = "{{.Recipe.CraftSection.Name}}"
  azs  = slice(data.aws_availability_zones.available.names, 0, 3)
  access_key = "{{.Recipe.Creds.AccessKey}}"
  secret_key = "{{.Recipe.Creds.SecretKey}}"
  region = "{{.Recipe.DefaultRegion}}"
}
`
