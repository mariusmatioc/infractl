package fargate

const LocalsNetwork = `
locals {
  name = "{{.Recipe.CraftSection.Name}}"
  azs  = slice(data.aws_availability_zones.available.names, 0, 3)
{{if .Recipe.Network.VpcId}}
  vpc = data.aws_vpc.vpc
{{else}}
  vpc = aws_vpc.vpc
{{end}}
}
`
