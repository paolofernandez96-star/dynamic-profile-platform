// ----------------------------------------------------
// FILE: iac/variables.tf
// Defines input variables for our root module.
// ----------------------------------------------------
variable "gcp_project_id" {
  description = "The GCP project ID to deploy resources into."
  type        = string
}

variable "gcp_region" {
  description = "The GCP region for resources."
  type        = string
  default     = "us-central1"
}


// ----------------------------------------------------
// FILE: iac/terraform.tfvars
// Assigns values to the variables.
// Create this file locally and do not commit to git.
// ----------------------------------------------------
# gcp_project_id = "your-gcp-project-id-here"