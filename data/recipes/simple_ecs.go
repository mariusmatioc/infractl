package recipes

const SIMPLE_ECS = `
craft:
  recipe_version: 0.1.0
  type: ecs
  docker_compose_files: []
  infra_env_file: ...	    # infrastructure related environment, such as AWS credentials
  network_craft_file: enter-name-here # if relative then to the location of this file

simple_ecs:
  type: fargate  # fargate or EC2
  # machine_type: t4g.micro  # required if type=EC2

  services:
    frontend:  # this should match the service name from docker-compose.yml
      desired_nodes: 2  # required
      cpu: 512  # required if type=fargate, 1024 = 1 vCPU unit
      memory: 4096  # required if type=fargate, in MiB units

      # domain_name: xxx.example.com  # optional, enables DNS record creation and certificate if load_balancer_https is specified
      # One of load_balancer_http or load_balancer_https is required
      # load_balancer_http: # optional
	  #  - 80:8080 # optional
      # load_balancer_https: # optional
	  #  - 443:8080 # optional
      # environment:  # optional
      #  BACKEND_HOST: externals.backend.backend # example
      # env_files: [".env"] #optional, if relative then to the parent of infractl root. In addition to the files in the docker compose

simple_rds:
  services:
    database:  # this should match the service name from docker-compose.yml
      engine: postgres  # required
      public: true  # defaults to false
      machine_type: db.t3.micro  # required
      storage_type: gp2 # default is old generation
      storage_iops: -1 # default is to disable specific iops speed
      storage_gigs: 20  # required
      env_files: [".env"] # Environment from here or docker compose must contain DB_USER, DB_PWD, DB_NAME ans DB_PORT

simple_mqs:
  services:
    queue:  # this should match the service name from docker-compose.yml
      engine: RabbitMQ  # required
      public: true  # defaults to false
      env_files: [".env"] # Environment from here or otherwise must contain RABBITMQ_PWD, RABBITMQ_PWD

externals: # optional
  - name: backend
    craft_file: backend.yml # if relative, then to the location of this file

credentials: # optional
  # Enter values here (optional). If missing, taken infra_env_file

  AWS_ACCESS_KEY_ID: ...
  AWS_SECRET_ACCESS_KEY: ...
  AWS_DEFAULT_REGION: ...
`
