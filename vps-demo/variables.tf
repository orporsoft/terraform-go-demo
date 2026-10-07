variable "container_name" {
  description = "Docker container name"
  type        = string
}

variable "external_port" {
  description = "Localhost port on VPS"
  type        = number
}

variable "image_tag" {
  description = "Docker image tag"
  type        = string
}

variable "app_env" {
  description = "Application environment"
  type        = string
}

variable "api_key" {
  description = "API key for the application"
  type        = string
  sensitive   = true
}

variable "postgres_password" {
  description = "PostgresSQL password"
  type        = string
  sensitive   = true
}