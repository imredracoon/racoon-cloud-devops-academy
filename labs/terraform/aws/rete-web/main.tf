# Lezioni 6.4 e 6.5: una VPC con due subnet pubbliche e alcuni server web.

# Data source: informazioni lette dal cloud, non create da Terraform.
data "aws_availability_zones" "disponibili" {
  state = "available"
}

# L'immagine Amazon Linux 2023 più recente, letta dal parametro pubblico di AWS.
data "aws_ssm_parameter" "al2023" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
}

locals {
  # Le prime due zone della regione, con un /24 ciascuna: { "eu-south-1a" = "10.20.0.0/24", ... }
  zone = slice(data.aws_availability_zones.disponibili.names, 0, 2)
  subnet_pubbliche = {
    for i, zona in local.zone : zona => cidrsubnet(var.cidr_vpc, 8, i)
  }
  prefisso = "${var.nome}-${var.ambiente}"
}

resource "aws_vpc" "principale" {
  cidr_block           = var.cidr_vpc
  enable_dns_hostnames = true

  tags = { Name = "${local.prefisso}-vpc" }
}

resource "aws_internet_gateway" "principale" {
  vpc_id = aws_vpc.principale.id

  tags = { Name = "${local.prefisso}-igw" }
}

resource "aws_subnet" "pubblica" {
  for_each = local.subnet_pubbliche

  vpc_id                  = aws_vpc.principale.id
  availability_zone       = each.key
  cidr_block              = each.value
  map_public_ip_on_launch = true

  tags = { Name = "${local.prefisso}-pubblica-${each.key}" }
}

resource "aws_route_table" "pubblica" {
  vpc_id = aws_vpc.principale.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.principale.id
  }

  tags = { Name = "${local.prefisso}-rt-pubblica" }
}

resource "aws_route_table_association" "pubblica" {
  for_each = aws_subnet.pubblica

  subnet_id      = each.value.id
  route_table_id = aws_route_table.pubblica.id
}

resource "aws_security_group" "web" {
  name        = "${local.prefisso}-web"
  description = "HTTP dal mondo"
  vpc_id      = aws_vpc.principale.id
}

resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.web.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "tcp"
  from_port         = 80
  to_port           = 80
}

resource "aws_vpc_security_group_egress_rule" "tutto" {
  security_group_id = aws_security_group.web.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

resource "aws_instance" "web" {
  count = var.numero_server

  ami                    = data.aws_ssm_parameter.al2023.value
  instance_type          = var.tipo_istanza
  subnet_id              = values(aws_subnet.pubblica)[count.index % length(local.zone)].id
  vpc_security_group_ids = [aws_security_group.web.id]

  user_data = templatefile("${path.module}/web.sh.tftpl", {
    nome_server = "${local.prefisso}-web-${count.index}"
  })
  user_data_replace_on_change = true

  metadata_options {
    http_tokens = "required" # solo IMDSv2
  }

  tags = { Name = "${local.prefisso}-web-${count.index}" }
}
