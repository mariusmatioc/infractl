package lambda

const Locals = `
locals {
  vpc = data.aws_vpc.vpc
  name = "{{.CraftSection.Name}}"
  access_key = "{{.Creds.AccessKey}}"
  secret_key = "{{.Creds.SecretKey}}"
  region = "{{.DefaultRegion}}"
}
`
