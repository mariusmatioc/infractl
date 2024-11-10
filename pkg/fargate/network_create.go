package fargate

const NetworkCreate = `

resource "aws_ecs_cluster" "main" {
	name = "{{.ClusterName}}"
	setting {
		name  = "containerInsights"
		value = "enabled"
	}
}

/*==== The VPC ======*/
{{if .VpcId}}
data "aws_vpc" "vpc" {
  id = "{{.VpcId}}"
}
{{else}}
resource "aws_vpc" "vpc" {
	cidr_block           = "{{.VpcCidr}}"
	enable_dns_hostnames = true
	enable_dns_support   = true
	tags = {
		Name        = "${local.name}-vpc"
	}
}
{{end}}

/* Internet gateway for the public subnet */
resource "aws_internet_gateway" "ig" {
	vpc_id = aws_vpc.vpc.id
	tags = {Name = "${local.name}-igw"	}
}

/* Elastic IP for NAT */
resource "aws_eip" "nat_eip" {
	domain        = "vpc"
	depends_on = [aws_internet_gateway.ig]
}

/* NAT */
resource "aws_nat_gateway" "nat" {
	allocation_id = "${aws_eip.nat_eip.id}"
	subnet_id     = "${element(aws_subnet.public_subnet.*.id, 0)}"
	depends_on    = [aws_internet_gateway.ig]
	tags = {Name        = "${local.name}-nat"	}
}

/*==== Subnets ======*/
/* Public subnets */
resource "aws_subnet" "public_subnet" {
	vpc_id                  = aws_vpc.vpc.id
	cidr_block              = cidrsubnet(local.vpc.cidr_block, 8, 0)
	map_public_ip_on_launch = true
    availability_zone = data.aws_availability_zones.available.names[0]
	tags = { Name = "${local.name}-public-subnet"}
}

resource "aws_subnet" "public_subnet_2" {
	vpc_id                  = aws_vpc.vpc.id
	cidr_block              = cidrsubnet(local.vpc.cidr_block, 8, 10)
	map_public_ip_on_launch = true
    availability_zone = data.aws_availability_zones.available.names[1]
	tags = {Name = "${local.name}-public-subnet-2"}
}

/* Private subnets */
resource "aws_subnet" "private_subnet" {
	vpc_id                  = aws_vpc.vpc.id
	cidr_block              = cidrsubnet(local.vpc.cidr_block, 8, 20)
	map_public_ip_on_launch = false
    availability_zone = data.aws_availability_zones.available.names[0]
	tags = {
		Name        = "${local.name}-private-subnet"
	}
}

resource "aws_subnet" "private_subnet_2" {
	vpc_id                  = aws_vpc.vpc.id
	cidr_block              = cidrsubnet(local.vpc.cidr_block, 8, 30)
    availability_zone = data.aws_availability_zones.available.names[1]
	tags = {
		Name        = "${local.name}-private-subnet-2"
	}
}

/* Routing table for private subnet */
resource "aws_route_table" "private" {
	vpc_id = aws_vpc.vpc.id
	tags = {Name        = "${local.name}-private-route-table"}
}

/* Routing table for public subnets */
resource "aws_route_table" "public" {
	vpc_id = aws_vpc.vpc.id
	tags = {
		Name        = "${local.name}-public-route-table"
	}
}

resource "aws_route" "public_internet_gateway" {
	route_table_id         = "${aws_route_table.public.id}"
	destination_cidr_block = "0.0.0.0/0"
	gateway_id             = "${aws_internet_gateway.ig.id}"
}

resource "aws_route" "private_nat_gateway" {
	route_table_id         = "${aws_route_table.private.id}"
	destination_cidr_block = "0.0.0.0/0"
	nat_gateway_id         = "${aws_nat_gateway.nat.id}"
}

/* Route table associations */
resource "aws_route_table_association" "public" {
	subnet_id      = "${aws_subnet.public_subnet.id}"
	route_table_id = "${aws_route_table.public.id}"
}

resource "aws_route_table_association" "private" {
	subnet_id      = "${aws_subnet.private_subnet.id}"
	route_table_id = "${aws_route_table.private.id}"
}

resource "aws_route_table_association" "public2" {
	subnet_id      = aws_subnet.public_subnet_2.id
	route_table_id = "${aws_route_table.public.id}"
}

resource "aws_route_table_association" "private2" {
	subnet_id      = "${aws_subnet.private_subnet_2.id}"
	route_table_id = "${aws_route_table.private.id}"
}

output "vpc_id" {
  value = aws_vpc.vpc.id
}

output "cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "public_subnet_id" {
  value = aws_subnet.public_subnet.id
}

output "public_subnet_2_id" {
  value = aws_subnet.public_subnet_2.id
}

output "private_subnet_id" {
  value = aws_subnet.private_subnet.id
}

output "private_subnet_2_id" {
  value = aws_subnet.private_subnet_2.id
}
`
