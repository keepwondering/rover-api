variable "project_id" { type = string }
variable "region" { type = string }
variable "service_name" { type = string }
variable "image" { type = string }
variable "domain_name" { type = string } # "" to skip mapping

# defaults so callers don't have to set these every time
variable "allow_unauthenticated" {
  type    = bool
  default = true
}
variable "min_instances" {
  type    = number
  default = 0
}
variable "max_instances" {
  type    = number
  default = 3
}
variable "cpu" {
  type    = string
  default = "1"
} # "1", "2", etc.
variable "memory" {
  type    = string
  default = "512Mi"
} # "512Mi", "1Gi", etc.

variable "env" {
  description = "Container env vars"
  type        = map(string)
  default     = {}
}
