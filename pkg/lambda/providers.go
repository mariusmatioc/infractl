package lambda

const Providers = `
terraform {
  required_version = ">= 1.5.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.37.0"
    }
  }
{{if ne .BackendConfig.AccessKey ""}}
  backend "s3" {
    bucket         = "infractl-terraform-backend"
    key            = "{{.BackendConfig.Organization}}/{{.BackendConfig.Project}}/{{.CraftSection.Name}}/terraform.tfstate"
    encrypt        = true
    dynamodb_table = "{{.BackendConfig.Organization}}-backend"

    region         = "us-west-2"
	access_key     = "{{.BackendConfig.AccessKey}}"
	secret_key     = "{{.BackendConfig.SecretKey}}"
  }
{{end}}
}

provider "aws" {
  region = "{{.DefaultRegion}}"
  access_key = "{{.Creds.AccessKey}}"
  secret_key = "{{.Creds.SecretKey}}"
}
`
