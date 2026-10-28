terraform {
  required_version = ">= 1.12"

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.9"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9"
    }
  }
}

resource "random_pet" "nome" {
  length = 2
}

resource "local_file" "benvenuto" {
  filename = "${path.module}/out/benvenuto.txt"
  content  = "Ciao! Questo server si chiama ${random_pet.nome.id}.\n"
}

output "nome_server" {
  value = random_pet.nome.id
}
