package fargate

const TaskDef = `
resource "aws_ecs_task_definition" "{{.Name}}" {
  family                   = "${local.name}-{{.Name}}-{{.TaskHash}}"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_execution_role.arn
  cpu                      = {{.Cpu}}
  memory                   = {{.Memory}}
  container_definitions = jsonencode([
    {
      name:        "${local.name}-{{.Name}}-{{.TaskHash}}",
      image:       {{.Image}},
      essential: true,
{{if .EnvsString}}
      environment: {{.EnvsString}},
{{end}}
{{if .Entrypoint}}
      entrypoint: {{.Entrypoint}},
{{end}}
      portMappings: [
{{range $val := .Ports}}
        {
			protocol:      "tcp",
			containerPort: {{$val.Target}},
			hostPort:      {{$val.Published}},
        }
{{end}}
      ],
{{if .HealthCheck.Test}}
	  healthCheck: {
		command:     {{.HealthCheck.Test}},
{{if .HealthCheck.Interval}} 
		interval:    {{.HealthCheck.Interval}},
{{end}}
{{if .HealthCheck.Timeout}} 
		timeout:     {{.HealthCheck.Timeout}},
{{end}}
{{if .HealthCheck.Retries}} 
		retries:     {{.HealthCheck.Retries}},
{{end}}
	  },	
{{end}}
      logConfiguration: {
        logDriver: "awslogs",
        options: {
          awslogs-group :        aws_cloudwatch_log_group.main.name,
          awslogs-stream-prefix: "{{.Name}}",
          awslogs-region:       local.region,
        }
      },
      linuxParameters: {
        initProcessEnabled: true
      }
    }
 ])

  {{.DependsOn}}
}
`
