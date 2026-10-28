variable "regione" {
  description = "Regione AWS"
  type        = string
  default     = "eu-south-1"
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

variable "cidr_vpc" {
  description = "Intervallo di indirizzi della VPC"
  type        = string
  default     = "10.20.0.0/16"
}

variable "numero_server" {
  description = "Quanti server web creare (uno per zona, a rotazione)"
  type        = number
  default     = 2

  validation {
    condition     = var.numero_server >= 1 && var.numero_server <= 4
    error_message = "Crea da 1 a 4 server."
  }
}

variable "tipo_istanza" {
  description = "Tipo di istanza EC2"
  type        = string
  default     = "t3.micro"
}
