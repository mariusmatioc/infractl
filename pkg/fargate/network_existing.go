package fargate

const NetworkExisting = `

data "aws_ecs_cluster" "main" {
	cluster_name  = "{{.ClusterName}}"
}

/*==== The VPC ======*/
data "aws_vpc" "vpc" {
  id = "{{.VpcId}}"
}

/*==== Subnets ======*/
/* Public subnets */
data "aws_subnet" "public_subnet" {
  filter {
    name   = "tag:Name"
    values = ["${local.name}-public-subnet"]
  }
}

data "aws_subnet" "public_subnet_2" {
  filter {
    name   = "tag:Name"
    values = ["${local.name}-public-subnet-2"]
  }
}

/* Private subnets */
data "aws_subnet" "private_subnet" {
  filter {
    name   = "tag:Name"
    values = ["${local.name}-private-subnet"]
  }
}

data "aws_subnet" "private_subnet_2" {
  filter {
    name   = "tag:Name"
    values = ["${local.name}-private-subnet-2"]
  }
}

output "vpc_id" {
  value = data.aws_vpc.vpc.id
}

output "cluster_name" {
  value = "{{.ClusterName}}"
}

output "public_subnet_id" {
  value = data.aws_subnet.public_subnet.id
}

output "public_subnet_2_id" {
  value = data.aws_subnet.public_subnet_2.id
}

output "private_subnet_id" {
  value = data.aws_subnet.private_subnet.id
}

output "private_subnet_2_id" {
  value = data.aws_subnet.private_subnet_2.id
}
`
