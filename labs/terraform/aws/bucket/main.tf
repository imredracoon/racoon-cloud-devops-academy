# Lezione 6.2: il primo bucket S3 gestito con Terraform.
terraform {
  required_version = ">= 1.12"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

provider "aws" {
  region = var.regione

  # Etichette applicate automaticamente a ogni risorsa che le supporta.
  default_tags {
    tags = {
      progetto = "academy"
      gestito  = "terraform"
    }
  }
}

variable "regione" {
  description = "Regione AWS in cui creare il bucket"
  type        = string
  default     = "eu-south-1"
}

# I nomi dei bucket sono unici in tutto il mondo: aggiungiamo un suffisso casuale.
resource "random_id" "suffisso" {
  byte_length = 4
}

resource "aws_s3_bucket" "dati" {
  bucket = "academy-dati-${random_id.suffisso.hex}"
}

resource "aws_s3_bucket_versioning" "dati" {
  bucket = aws_s3_bucket.dati.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_public_access_block" "dati" {
  bucket = aws_s3_bucket.dati.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

output "nome_bucket" {
  value = aws_s3_bucket.dati.bucket
}
