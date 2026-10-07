output "container_name" {
  value = docker_container.app.name
}

output "local_url" {
  value = "http://127.0.0.1:${var.external_port}"
}
