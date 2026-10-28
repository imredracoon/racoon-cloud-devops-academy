variable "progetto" {
  description = "ID del progetto Google Cloud"
  type        = string
}

variable "regione" {
  description = "Regione Google Cloud"
  type        = string
  default     = "europe-west8"
}

variable "nome" {
  description = "Prefisso dei nomi delle risorse"
  type        = string
  default     = "academy"
}

variable "ambiente" {
  description = "Ambiente: dev, test o prod"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "test", "prod"], var.ambiente)
    error_message = "L'ambiente deve essere dev, test o prod."
  }
}

variable "cidr_subnet" {
  description = "Intervallo di indirizzi della subnet"
  type        = string
  default     = "10.20.0.0/24"
}

variable "numero_server" {
  description = "Quanti server web creare (distribuiti sulle zone)"
  type        = number
  default     = 2

  validation {
    condition     = var.numero_server >= 1 && var.numero_server <= 4
    error_message = "Crea da 1 a 4 server."
  }
}

variable "tipo_macchina" {
  description = "Tipo di macchina Compute Engine"
  type        = string
  default     = "e2-micro"
}
