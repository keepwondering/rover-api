output "service_url" {
  value = google_cloud_run_v2_service.svc.uri
}
output "domain_name" {
  value = var.domain_name
}
