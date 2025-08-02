# iac/main.tf

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 6.0"
    }
  }
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

provider "google-beta" {
  project = var.gcp_project_id
  region  = var.gcp_region
}


# --- Database Module ---
# Creates our Cloud SQL and Redis instances.
module "database" {
  source           = "./modules/gcp_database"
  project_id       = var.gcp_project_id
  region           = var.gcp_region
  db_instance_name = "profile-db-main"
  redis_name       = "profile-cache-main"
}

# --- Services Module ---
# Creates our backend microservices on Cloud Run.
module "services" {
  source          = "./modules/gcp_services"
  project_id      = var.gcp_project_id
  region          = var.gcp_region
  service_name    = "profile-service"
  container_image = "gcr.io/${var.gcp_project_id}/profile-service:latest"
}

# --- API Gateway Resources ---
resource "google_api_gateway_api" "profile_api" {
  provider = google-beta
  project  = var.gcp_project_id
  api_id   = "profile-api"
}

resource "google_api_gateway_api_config" "profile_api_config" {
  provider      = google-beta
  project       = var.gcp_project_id
  api           = google_api_gateway_api.profile_api.api_id
  api_config_id = "profile-api-config"

  openapi_documents {
    document {
      path     = "openapi.yaml"
      contents = filebase64("${path.module}/openapi.yaml")
    }
  }
  
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "profile_gateway" {
  provider    = google-beta
  project     = var.gcp_project_id
  region      = var.gcp_region
  gateway_id  = "profile-gateway"
  api_config  = google_api_gateway_api_config.profile_api_config.id
}

output "api_gateway_url" {
  description = "The default URL of the profile gateway"
  value       = google_api_gateway_gateway.profile_gateway.default_hostname
}