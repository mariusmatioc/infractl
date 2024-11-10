package fargate

const Service = `
{{ $name := .Name }}
{{ $taskHash := .TaskHash }}
resource "aws_ecs_service" "{{.Name}}" {
  name                               = "${local.name}-{{.Name}}"
  cluster                            = data.aws_ecs_cluster.main.id
  task_definition                    = aws_ecs_task_definition.{{.Name}}.arn
  desired_count                      = {{.DesiredCount}}
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"
  enable_execute_command             = true
  network_configuration {
    security_groups  = [aws_security_group.{{.Name}}.id]
    subnets          = [data.aws_subnet.private_subnet.id]
    assign_public_ip = false
  }
{{range $p := .LoadBalacerTargets}}
  load_balancer {
    target_group_arn = aws_lb_target_group.{{$name}}-{{$p}}.arn
    container_name   = "${local.name}-{{$name}}-{{$taskHash}}"
    container_port   = {{$p}}
  }
{{end}}
  lifecycle { ignore_changes = [desired_count] }
}

resource "aws_lb" "{{.Name}}" {
  name               = "${local.name}-{{.Name}}"
  internal           = false
  load_balancer_type = "network"
  subnets            = [data.aws_subnet.public_subnet.id, data.aws_subnet.public_subnet_2.id]
  enable_cross_zone_load_balancing = true
  enable_deletion_protection = false
}

{{range $lbports := .LoadBalancerPortsTCP}}
resource "aws_lb_listener" "{{$name}}-{{$lbports.Published}}" {
  load_balancer_arn = aws_lb.{{$name}}.arn
  port              = {{$lbports.Published}}
  protocol          = "TCP"

  default_action {
    target_group_arn = aws_lb_target_group.{{$name}}-{{$lbports.Target}}.id
    type             = "forward"
  }
}
{{end}}

{{range $lbports := .LoadBalancerPortsTLS}}
resource "aws_lb_listener" "{{$name}}-{{$lbports.Published}}" {
  load_balancer_arn = aws_lb.{{$name}}.arn
  port              = {{$lbports.Published}}
  protocol          = "TLS"
  certificate_arn   = aws_acm_certificate_validation.{{$name}}.certificate_arn

  default_action {
    target_group_arn = aws_lb_target_group.{{$name}}-{{$lbports.Target}}.id
    type             = "forward"
  }
}
{{end}}

{{range $p := .LoadBalacerTargets}}
resource "aws_lb_target_group" "{{$name}}-{{$p}}" {
  port        = {{$p}}
  protocol    = "TCP"
  vpc_id      = data.aws_vpc.vpc.id
  target_type = "ip"
}
{{end}}

{{if .DomainName}}
data "aws_route53_zone" "{{.Name}}" {
  name         = "{{.HostedZone}}"
  private_zone = false
}

resource "aws_route53_record" "{{.Name}}" {
  zone_id = data.aws_route53_zone.{{.Name}}.zone_id
  name    = "{{.DomainName}}"
  type    = "A"

  alias {
    name                   = aws_lb.{{.Name}}.dns_name
    zone_id                = aws_lb.{{.Name}}.zone_id
    evaluate_target_health = true
  }
}

{{if .NeedsCertificate}}
resource "aws_acm_certificate" "{{.Name}}" {
  domain_name                 = aws_route53_record.{{.Name}}.name
  validation_method           = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "{{.Name}}-cert" {
  for_each = {
    for dvo in aws_acm_certificate.{{.Name}}.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.{{.Name}}.zone_id
}

resource "aws_acm_certificate_validation" "{{.Name}}" {
  certificate_arn         = aws_acm_certificate.{{.Name}}.arn
  validation_record_fqdns = [for record in aws_route53_record.{{.Name}}-cert : record.fqdn]
}
{{end}}
{{end}}

output "{{.Name}}" {
{{if .DomainName}}
  value = aws_route53_record.{{.Name}}.name
{{else}}
  value = aws_lb.{{.Name}}.dns_name
{{end}}
}

`
