variable "container_name" {
  description = "Name of the container"
  type = string
}
variable "container_version" {
  description = "Version of the Docker image"
  type        = string
  default     = "latest-dev"
}

variable "env" {
  description = "Environment to deploy (prod or dev)"
  type        = string
  default     = "dev"
}

variable "db_username" {
  description = "Database username"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "db_ip" {
  description = "Database IP address"
  type        = string
  sensitive   = true
}

variable "db_port" {
  description = "Database port"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
  sensitive   = true
}

variable "api_secret" {
  description = "API secret"
  type        = string
  sensitive   = true
}

variable "graylog_host" {
  description = "The Graylog host with ip and port: <ip>:<port>"
  type        = string
}