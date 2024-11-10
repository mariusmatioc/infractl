package fargate

const NetworkUse = `

data "aws_ecs_cluster" "main" {
	cluster_name = "{{.ClusterName}}"
}

data "aws_vpc" "vpc" {
  	id = "{{.VpcId}}"
}

data "aws_subnet" "public_subnet" {
	id = "{{.PublicSubnetId}}"
}

data "aws_subnet" "public_subnet_2" {
	id = "{{.PublicSubnet2Id}}"
}

data "aws_subnet" "private_subnet" {
	id = "{{.PrivateSubnetId}}"
}

data "aws_subnet" "private_subnet_2" {
	id = "{{.PrivateSubnet2Id}}"
}
`
