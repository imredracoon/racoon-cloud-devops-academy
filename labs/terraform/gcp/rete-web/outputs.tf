output "rete" {
  description = "Nome della rete VPC"
  value       = google_compute_network.principale.name
}

output "indirizzi_web" {
  description = "URL dei server web"
  value       = [for i in google_compute_instance.web : "http://${i.network_interface[0].access_config[0].nat_ip}"]
}
