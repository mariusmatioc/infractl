package fargate

const Providers = `
terraform {
  required_version = ">= 1.5.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.7.0"
    }

    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0.2"
    }
  }
{{if ne .BackendConfig.Bucket ""}}
  backend "s3" {
    bucket         = "{{.BackendConfig.Bucket}}"
    key            = "{{.BackendConfig.Key}}/{{.CraftSection.Name}}/terraform.tfstate"
    encrypt        = true
    dynamodb_table = "{{.BackendConfig.Bucket}}"
  }
{{end}}
}

provider "aws" {
  region = "{{.DefaultRegion}}"
  access_key = "{{.Creds.AccessKey}}"
  secret_key = "{{.Creds.SecretKey}}"
}

provider "docker" {
}
`
