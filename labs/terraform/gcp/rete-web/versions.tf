terraform {
  required_version = ">= 1.12"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 8.0"
    }
  }
}

provider "google" {
  project = var.progetto
  region  = var.regione

  default_labels = {
    progetto = var.nome
    ambiente = var.ambiente
    gestito  = "terraform"
  }
}
