package fargate

const ServiceDiscovery = `
resource "aws_service_discovery_service" "{{.Name}}" {
  name = "${local.name}-{{.Name}}"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.cloudmap.id

    dns_records {
      ttl  = 60
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}
`
