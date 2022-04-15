
variable "lb_target_group_names" {
    type = list(string)
}

variable "lb_target_type" {
    type = string
}
variable "target_ids_arn" {
    type = list(string)
}
variable "load_balancer_name" {
    type = string
}

variable "development_environment" {
    type = string
}

variable "security_groups" {
    type = list(string)
}

variable "subnet_ids" {
    type = list(string)
}