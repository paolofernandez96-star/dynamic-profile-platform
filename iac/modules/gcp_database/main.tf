// ----------------------------------------------------
// FILE: iac/modules/gcp_database/main.tf
// This module is responsible for creating database resources.
// ----------------------------------------------------
resource "google_sql_database_instance" "postgres" {
  name             = var.db_instance_name
  project          = var.project_id
  region           = var.region
  database_version = "POSTGRES_14"
  
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled = true
    }
  }
}

resource "google_redis_instance" "cache" {
  name           = var.redis_name
  project        = var.project_id
  region         = var.region
  tier           = "BASIC"
  memory_size_gb = 1
}