terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~>3.0"
    }
  }
}

provider "docker" {
  host = "ssh://terraform@orporsoft.com"

  registry_auth {
    address     = "ghcr.io"
    config_file = pathexpand("~/.docker/config.json")
  }
}

resource "docker_network" "app_network" {
  name = "terraform-go-network"
}

resource "docker_image" "app" {
  name = "ghcr.io/orporsoft/terraform-go-demo:${var.image_tag}"
}

resource "docker_container" "app" {
  name  = var.container_name
  image = docker_image.app.image_id

  env = [
    "APP_ENV=${var.app_env}",
    "API_KEY=${var.api_key}",
    "DB_HOST=postgres",
    "DB_PORT=5432",
    "DB_NAME=orderdb",
    "DB_USER=orderuser",
    "DB_PASSWORD=${var.postgres_password}"
  ]

  depends_on = [
    docker_container.postgres
  ]

  ports {
    internal = 8080
    external = var.external_port
    ip       = "127.0.0.1"
  }

  networks_advanced {
    name = docker_network.app_network.name
  }
}

resource "docker_image" "postgres" {
  name = "postgres:16-alpine"
}

resource "docker_volume" "postgres_data" {
  name = "terraform-go-postgres-data"
}

resource "docker_container" "postgres" {
  name  = "terraform-go-postgres"
  image = docker_image.postgres.image_id

  env = [
    "POSTGRES_DB=orderdb",
    "POSTGRES_USER=orderuser",
    "POSTGRES_PASSWORD=${var.postgres_password}"
  ]

  networks_advanced {
    name = docker_network.app_network.name
  }

  volumes {
    volume_name    = docker_volume.postgres_data.name
    container_path = "/var/lib/postgresql/data"
  }
}

