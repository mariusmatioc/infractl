package fargate

// AWS managed message queue (only RabbitMQ for now)
const Mq = `
resource "aws_mq_broker" "{{.Name}}" {
  broker_name = "${local.name}-{{.Name}}"

  engine_type        = "RabbitMQ"
  engine_version     = "3.11.20"
  host_instance_type = "mq.m5.large"
  publicly_accessible = {{.Public}}
{{if not .Public}}
  subnet_ids = [data.aws_subnet.private_subnet.id]
  security_groups    = [aws_security_group.{{.Name}}.id]
{{end}}
  user {
    username = "{{.UserName}}"
    password = "{{.Password}}"
  }
}

{{if not .Public}}
resource "aws_security_group" "{{.Name}}" {
  name   = "${local.name}-{{.Name}}"
  vpc_id = data.aws_vpc.vpc.id

  ingress {
    protocol         = "tcp"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = [data.aws_vpc.vpc.cidr_block]
  }

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}
{{end}}

output "{{.Name}}" {
  value = aws_mq_broker.{{.Name}}.instances.0.endpoints.0
}

`
