resource "google_cloud_run_v2_service" "svc" {
  name     = var.service_name
  location = var.region

  template {
    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }

    containers {
      image = var.image
      resources {
        cpu_idle = true
        limits = {
          "cpu"    = var.cpu
          "memory" = var.memory
        }
      }

      dynamic "env" {
        for_each = var.env
        content {
          name  = env.key
          value = env.value
        }
      }
      # Cloud Run expects the app to listen on $PORT (default 8080)
    }

    max_instance_request_concurrency = 80
  }
}

resource "google_cloud_run_v2_service_iam_member" "public" {
  count    = var.allow_unauthenticated ? 1 : 0
  name     = google_cloud_run_v2_service.svc.name
  location = var.region
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_domain_mapping" "mapping" {
  provider = google-beta
  count    = var.domain_name != "" ? 1 : 0

  location = var.region
  name     = var.domain_name

  metadata { namespace = var.project_id }

  spec {
    route_name = google_cloud_run_v2_service.svc.name
  }

  depends_on = [google_cloud_run_v2_service.svc]
}
