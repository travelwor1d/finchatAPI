locals {
  project = var.project
  region  = "us-west1"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.64.0"
    }
  }
  backend "gcs" {
    bucket = "finchat-api-terraform-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = local.project
  region  = local.region
}

data "google_project" "project" {
}

resource "google_project_service" "iam" {
  project = local.project
  service = "iam.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "cloudbuild" {
  project = local.project
  service = "cloudbuild.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "run" {
  project = local.project
  service = "run.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "cloudidentity" {
  project = local.project
  service = "cloudidentity.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "sqladmin" {
  project = local.project
  service = "sqladmin.googleapis.com"

  disable_dependent_services = true
}

resource "google_sql_database_instance" "db" {
  name             = "finchat-db-staging"
  database_version = "MYSQL_8_0"
  region           = local.region

  settings {
    tier            = "db-f1-micro"
    disk_autoresize = true
  }

  deletion_protection = true
}

resource "google_sql_user" "default" {
  name     = var.db_username
  instance = google_sql_database_instance.db.name
  password = var.db_password
}

resource "google_sql_database" "core" {
  name     = "core"
  instance = google_sql_database_instance.db.name
}

locals {
  db_connection_string = "${google_sql_user.default.name}:${google_sql_user.default.password}@unix(/cloudsql/${google_sql_database_instance.db.connection_name})/${google_sql_database.core.name}?parseTime=true"
}

data "google_container_registry_image" "finchat_api" {
  name = "finchat-api"
}

resource "google_cloud_run_service" "api" {
  name     = "finchat-api-staging"
  location = local.region

  template {
    spec {
      containers {
        image = "gcr.io/${local.project}/finchat-api:${var.image_tag}"
        ports {
          container_port = 8080
        }
        env {
          name  = "MYSQL_CONN_STR"
          value = local.db_connection_string
        }
        env {
          name  = "TWILIO_SID"
          value = var.twilio["sid"]
        }
        env {
          name  = "TWILIO_TOKEN"
          value = var.twilio["token"]
        }
        env {
          name  = "TWILIO_VERIFY"
          value = var.twilio["verify"]
        }
        env {
          name  = "STRIPE_KEY"
          value = var.stripe["key"]
        }
        env {
          name  = "PUB_KEY"
          value = var.pubnub["pubKey"]
        }
        env {
          name  = "SUB_KEY"
          value = var.pubnub["subKey"]

        }
        env {
          name  = "SEC_KEY"
          value = var.pubnub["secKey"]
        }
        env {
          name  = "SERVER_UUID"
          value = "finchat-staging-server"
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "100"
        "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.db.connection_name
      }
    }
  }

  depends_on = [
    google_project_service.run,
    google_project_service.sqladmin
  ]
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.api.location
  project  = google_cloud_run_service.api.project
  service  = google_cloud_run_service.api.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
