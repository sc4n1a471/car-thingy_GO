terraform {
  required_providers {
    docker = {
      source = "registry.opentofu.org/kreuzwerker/docker"
      version = "~> 3.9.0"
    }
  }
}

locals {
  container_name = var.container_name
}

resource "docker_image" "car-thingy_go" {
  name = "sc4n1a471/car-thingy_go:${var.container_version}"
}

resource "docker_container" "car-thingy_go" {
  name  = local.container_name
  image = docker_image.car-thingy_go.name

  volumes {
    host_path      = var.env == "prod" ? "/media/car-thingy/prod" : "/media/car-thingy/dev"
    container_path = "/app/logs"
  }

  volumes {
    host_path      = var.env == "prod" ? "/media/car-thingy/prod" : "/media/car-thingy/dev"
    container_path = "/app/downloaded_images"
  }

  env = [
    "DB_USERNAME=${var.db_username}",
    "DB_PASSWORD=${var.db_password}",
    "DB_IP=${var.db_ip}",
    "DB_PORT=${var.db_port}",
    "DB_NAME=${var.db_name}",
    "API_SECRET=${var.api_secret}",
  ]

  ports {
    internal = 3000
    external = var.env == "prod" ? 3010 : (var.env == "dev" ? 3011 : null)
  }

  log_driver = "gelf"
  log_opts = {
    "gelf-address" = "udp://${var.graylog_host}"
  }

  networks_advanced {
    name = "car-thingy"
    ipv4_address = var.env == "prod" ? "172.18.0.5" : (var.env == "dev" ? "172.18.0.4" : null)
  }

  restart = "on-failure"
  max_retry_count = 5
}