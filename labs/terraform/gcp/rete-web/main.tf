# Lezioni 6.4 e 6.5: una rete VPC con una subnet e alcuni server web.

# Data source: informazioni lette dal cloud, non create da Terraform.
data "google_compute_zones" "disponibili" {
  region = var.regione
}

data "google_compute_image" "debian" {
  family  = "debian-12"
  project = "debian-cloud"
}

locals {
  zone     = slice(data.google_compute_zones.disponibili.names, 0, 2)
  prefisso = "${var.nome}-${var.ambiente}"
}

resource "google_compute_network" "principale" {
  name                    = "${local.prefisso}-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "web" {
  name          = "${local.prefisso}-web"
  network       = google_compute_network.principale.id
  region        = var.regione
  ip_cidr_range = var.cidr_subnet
}

resource "google_compute_firewall" "http" {
  name    = "${local.prefisso}-consenti-http"
  network = google_compute_network.principale.id

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web"]
}

resource "google_compute_instance" "web" {
  count = var.numero_server

  name         = "${local.prefisso}-web-${count.index}"
  zone         = local.zone[count.index % length(local.zone)]
  machine_type = var.tipo_macchina
  tags         = ["web"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.debian.self_link
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.web.id

    # Un indirizzo IP esterno temporaneo.
    access_config {}
  }

  metadata_startup_script = templatefile("${path.module}/web.sh.tftpl", {
    nome_server = "${local.prefisso}-web-${count.index}"
  })

  shielded_instance_config {
    enable_secure_boot = true
  }
}
