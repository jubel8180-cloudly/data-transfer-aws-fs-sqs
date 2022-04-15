
output "vpc_id" {
  value = "${aws_vpc.head.id}"
}

output "subnet_ids" {
  value = [for data in aws_subnet.head: data.id]
}

output "security_group"{
    value = aws_security_group.public.id
}