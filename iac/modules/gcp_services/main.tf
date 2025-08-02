// ----------------------------------------------------
// FILE: iac/modules/gcp_services/main.tf
// This module is responsible for creating Cloud Run services.
// ----------------------------------------------------
resource "google_cloud_run_v2_service" "profile_service" {
  name     = var.service_name
  location = var.region
  project  = var.project_id

  deletion_protection = false

  template {
    # THIS IS THE CORRECTED ANNOTATION (no hyphen in cloudsql-instances)
    annotations = {
      "run.googleapis.com/cloudsql-instances" = "dynamic-profile-platform:us-central1:profile-db-main"
    }

    containers {
      image = var.container_image
      
      env {
        name  = "DB_USER"
        value = "postgres"
      }
      env {
        name  = "DB_NAME"
        value = "profiles_db"
      }
      env {
        name  = "INSTANCE_CONNECTION_NAME"
        value = "dynamic-profile-platform:us-central1:profile-db-main"
      }
      env {
        name = "DB_PASS"
        value_source {
          secret_key_ref {
            secret  = "postgres-password"
            version = "latest"
          }
        }
      }
    }
  }
}