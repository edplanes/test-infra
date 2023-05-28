
# Authorization 
variable "hcloud_token" {
  sensitive = true

  description = "Hetzner Cloud API token"

  validation {
    condition     = length(var.hcloud_token) == 64
    error_message = "Please provide a valid Hetzner Cloud API token"
  }
}

variable "cloudflare_api_token" {
  sensitive = true

  description = "Cloudflare API token"

  validation {
    condition     = length(var.cloudflare_api_token) > 0
    error_message = "Please provide a valid Cloudflare API token"
  }
}

# Hetzner Cloud

variable "hcloud_ssh_keys" {
  default = []

  description = "List of SSH key IDs to inject into the server"
}

variable "main_server_count" {
  default = 1

  description = "Number of main nodes to create"
}

variable "agents_server_count" {
  default = 0

  description = "Number of agents nodes to create"
}

variable "data_server_count" {
  default = 1

  description = "Number of data nodes to create"
}

variable "main_server_type" {
  default = "cx21"

  description = "Server type for main nodes"
}

variable "agents_server_type" {
  default = "cx21"

  description = "Server type for agents nodes"
}

variable "data_server_type" {
  default = "cx21"

  description = "Server type for data nodes"
}

# Cloudflare

variable "cloudflare_account_id" {
  sensitive = true

  description = "Cloudflare account ID"

  validation {
    condition     = length(var.cloudflare_account_id) > 0
    error_message = "Please provide a valid Cloudflare account ID"
  }
}