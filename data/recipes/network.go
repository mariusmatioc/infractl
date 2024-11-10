package recipes

const NETWORK = `
craft:
  recipe_version: 0.1.0
  type: network
  infra_env_file: ...	    # infrastructure related environment, such as AWS credentials

network:
  cluster_name: enter-name-here
  # One of vpc_id or vpc_cidr is required
  vpc_cidr: 10.0.0.0/16
  #vpc_id: vpc-0bc0ec9e9ddbc0ef7 # Uses existing VPC if it exists

credentials: # optional
  # Enter values here (optional). If missing, taken infra_env_file
  AWS_ACCESS_KEY_ID: ...
  AWS_SECRET_ACCESS_KEY: ...
  AWS_DEFAULT_REGION: ...
`
