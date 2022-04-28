variable "org_slug" {
  type        = string
  description = "Human-readable Org name/slug. The `citrix` in https://podio.com/citrix/hello-world"
}

variable "username" {
  type        = string
  sensitive   = true
  description = "Your Podio username"
}

variable "password" {
  type        = string
  sensitive   = true
  description = "Your Podio password"
}

variable "client_id" {
  type        = string
  sensitive   = true
  description = "Your API key ID/name"
}

variable "client_secret" {
  type        = string
  sensitive   = true
  description = "Your API key secret"
}
