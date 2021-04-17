variable "project" {
  type = string
}

variable "image_tag" {
  type        = string
  description = "Docker image tag of finchat-api"
  default     = "latest"
}

variable "db_username" {
  type      = string
  sensitive = true
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "twilio" {
  type      = map(string)
  sensitive = true
}

variable "stripe" {
  type      = map(string)
  sensitive = true
}

variable "pubnub" {
  type      = map(string)
  sensitive = true
}
