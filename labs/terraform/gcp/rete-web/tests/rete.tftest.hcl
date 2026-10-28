# Test senza cloud: il provider è simulato (mock), quindi non servono credenziali e non si crea nulla.
# Esecuzione: terraform test
mock_provider "google" {
  mock_data "google_compute_zones" {
    defaults = {
      names = ["europe-west8-a", "europe-west8-b", "europe-west8-c"]
    }
  }
}

variables {
  progetto = "progetto-di-prova"
}

run "server_distribuiti_sulle_zone" {
  command = plan

  variables {
    numero_server = 3
  }

  assert {
    condition     = google_compute_instance.web[1].zone == "europe-west8-b"
    error_message = "Il secondo server deve stare nella seconda zona."
  }

  assert {
    condition     = google_compute_instance.web[2].zone == "europe-west8-a"
    error_message = "Il terzo server deve tornare nella prima zona."
  }
}

run "ambiente_non_valido" {
  command = plan

  variables {
    ambiente = "collaudo"
  }

  expect_failures = [var.ambiente]
}
