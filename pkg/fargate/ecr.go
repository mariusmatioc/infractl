package fargate

const Ecr = `

resource "aws_ecr_repository" "{{.Name}}" {
  name                 = "${local.name}-{{.Name}}"
  image_tag_mutability = "MUTABLE"
  force_delete = true
}

resource "aws_ecr_repository_policy" "{{.Name}}-repo-policy" {
  repository = aws_ecr_repository.{{.Name}}.name
  policy     = <<EOF
  {
    "Version": "2008-10-17",
    "Statement": [
      {
        "Sid": "adds full ecr access to the repository",
        "Effect": "Allow",
        "Principal": "*",
        "Action": [
          "ecr:BatchCheckLayerAvailability",
          "ecr:BatchGetImage",
          "ecr:CompleteLayerUpload",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetLifecyclePolicy",
          "ecr:InitiateLayerUpload",
          "ecr:PutImage",
          "ecr:UploadLayerPart"
        ]
      }
    ]
  }
  EOF
}

resource "docker_image" "{{.Name}}_image" {
  name = "${local.name}-{{.Name}}"
  build {
    context = "{{.BuildContext}}"
{{if .BuildDockerfile}}
    dockerfile = "{{.BuildDockerfile}}"
{{end}}
{{if .BuildTarget}}
    target = "{{.BuildTarget}}"
{{end}}
    tag = ["${aws_ecr_repository.{{.Name}}.repository_url}", "{{.BuildHash}}"]
  }
  depends_on = [aws_ecr_repository.{{.Name}}]
  provisioner "local-exec" {
    environment = {
      AWS_ACCESS_KEY_ID = local.access_key
      AWS_SECRET_ACCESS_KEY = local.secret_key
    }
    command = "aws ecr get-login-password --region ${local.region} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${local.region}.amazonaws.com && docker push ${aws_ecr_repository.{{.Name}}.repository_url}"
  }
}

`
