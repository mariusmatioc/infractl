package recipes

const SIMPLE_LAMBDA = `
craft:
  recipe_version: 0.1.0
  type: lambda
  infra_env_file: ...	    # infrastructure related environment, such as AWS credentials
  network_craft_file: enter-name-here # if relative then to the location of this file

simple_lambda:
  function_name: ...
  handler: your_handler_function_name, should have signature func(event, context)
  runtime: python3.8 # to match your code, can be any runtime supported by AWS
  source_folder: ... # if relative then to the parent of infractl root
  memory_size: 1024  # MiB units
  timeout: 30  # seconds
  ephemeral_storage: 1000 # Optional MiB units
  layers: [] # optional, ARNs of Lambda layers
  triggers: # At least one of triggers is required
	schedule_expression: rate(1 day) # or cron(0 20 * * ? *) for 8pm UTC 
	s3_object_created: becket-name
  # environment:  # optional
  #  PARAMETER_KEY: PARAMETER_VALUE
  # env_files: [".env"] # optional, if relative then to the parent of infractl root

credentials: # optional
  # Enter values here (optional). If missing, taken infra_env_file
  AWS_ACCESS_KEY_ID: ...
  AWS_SECRET_ACCESS_KEY: ...
  AWS_DEFAULT_REGION: ...
`
