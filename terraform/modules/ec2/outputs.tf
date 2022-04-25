

output "aws_lb_target_group_arn" {
    value = aws_lb_target_group.main.arn
}

output "aws_lb_dns_name" {
    value = aws_lb.main.dns_name
}

