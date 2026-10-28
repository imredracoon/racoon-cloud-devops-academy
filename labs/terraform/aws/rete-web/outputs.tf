output "vpc_id" {
  description = "ID della VPC"
  value       = aws_vpc.principale.id
}

output "subnet_pubbliche" {
  description = "Mappa zona => ID della subnet"
  value       = { for zona, s in aws_subnet.pubblica : zona => s.id }
}

output "indirizzi_web" {
  description = "URL dei server web"
  value       = [for i in aws_instance.web : "http://${i.public_ip}"]
}
