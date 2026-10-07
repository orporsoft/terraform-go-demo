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

