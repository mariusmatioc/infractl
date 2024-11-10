package fargate

const SecurityGroup = `

resource "aws_security_group" "{{.Name}}" {
  name   = "${local.name}-{{.Name}}"
  vpc_id = data.aws_vpc.vpc.id

{{range $val := .Ports}}
  ingress {
    protocol         = "tcp"
    from_port        = {{$val.Target}}
    to_port          = {{$val.Target}}
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
{{end}}

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

}
`
