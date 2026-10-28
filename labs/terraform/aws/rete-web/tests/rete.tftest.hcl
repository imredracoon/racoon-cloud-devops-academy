# Test senza cloud: il provider è simulato (mock), quindi non servono credenziali e non si crea nulla.
# Esecuzione: terraform test
mock_provider "aws" {
  mock_data "aws_availability_zones" {
    defaults = {
      names = ["eu-south-1a", "eu-south-1b", "eu-south-1c"]
    }
  }

  mock_data "aws_ssm_parameter" {
    defaults = {
      value = "ami-0123456789abcdef0"
    }
  }
}

run "due_subnet_in_due_zone" {
  command = plan

  assert {
    condition     = length(aws_subnet.pubblica) == 2
    error_message = "Servono esattamente due subnet pubbliche."
  }

  assert {
    condition     = aws_subnet.pubblica["eu-south-1b"].cidr_block == "10.20.1.0/24"
    error_message = "La seconda subnet deve essere 10.20.1.0/24."
  }
}

run "server_distribuiti_sulle_zone" {
  command = plan

  variables {
    numero_server = 3
  }

  assert {
    condition     = length(aws_instance.web) == 3
    error_message = "Devono esserci tre server."
  }
}

run "ambiente_non_valido" {
  command = plan

  variables {
    ambiente = "collaudo"
  }

  expect_failures = [var.ambiente]
}
