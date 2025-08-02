// ----------------------------------------------------
// FILE: iac/modules/gcp_database/variables.tf
// ----------------------------------------------------
variable "project_id" { type = string }
variable "region" { type = string }
variable "db_instance_name" { type = string }
variable "redis_name" { type = string }