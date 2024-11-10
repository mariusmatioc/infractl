package fargate

const Rd = `

resource "aws_db_instance" "{{.Name}}" {
  identifier = "${local.name}-{{.DbName}}"

  engine            = "{{.DbEngine}}"
  allocated_storage = {{.StorageGigs}}
  instance_class    = "{{.MachineType}}"

  db_name           = "{{.DbName}}"
  username          = "{{.UserName}}"
  password          = "{{.Password}}"
  port              = {{.Port}}
  publicly_accessible = {{.Public}}

  vpc_security_group_ids = [aws_security_group.{{.Name}}.id]
  db_subnet_group_name   = aws_db_subnet_group.{{.Name}}.name

  multi_az            = false

  allow_major_version_upgrade = true
  auto_minor_version_upgrade  = true
  apply_immediately           = true

  deletion_protection      = false
  delete_automated_backups = true
  skip_final_snapshot = true

}

resource "aws_security_group" "{{.Name}}" {
  name   = "${local.name}-{{.Name}}"
  vpc_id = data.aws_vpc.vpc.id

  ingress {
    protocol         = "tcp"
    from_port        = {{.Port}}
    to_port          = {{.Port}}
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

}

resource "aws_db_subnet_group" "{{.Name}}" {
  name       = "${local.name}-{{.Name}}"
{{if .Public}}
  subnet_ids = [data.aws_subnet.public_subnet.id, data.aws_subnet.public_subnet_2.id]
{{else}}
  subnet_ids = [data.aws_subnet.private_subnet.id, data.aws_subnet.private_subnet_2.id]
{{end}}
}

output "{{.Name}}" {
  value = aws_db_instance.{{.Name}}.address
}

`
