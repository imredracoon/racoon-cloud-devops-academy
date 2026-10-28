# Lezione 6.2: il primo bucket Cloud Storage gestito con Terraform.
terraform {
  required_version = ">= 1.12"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 8.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

provider "google" {
  project = var.progetto
  region  = var.regione

  # Etichette applicate automaticamente a ogni risorsa che le supporta.
  default_labels = {
    progetto = "academy"
    gestito  = "terraform"
  }
}

variable "progetto" {
  description = "ID del progetto Google Cloud"
  type        = string
}

variable "regione" {
  description = "Regione Google Cloud in cui creare il bucket"
  type        = string
  default     = "europe-west8"
}

# I nomi dei bucket sono unici in tutto il mondo: aggiungiamo un suffisso casuale.
resource "random_id" "suffisso" {
  byte_length = 4
}

resource "google_storage_bucket" "dati" {
  name     = "academy-dati-${random_id.suffisso.hex}"
  location = var.regione

  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  versioning {
    enabled = true
  }
}

output "nome_bucket" {
  value = google_storage_bucket.dati.name
}
