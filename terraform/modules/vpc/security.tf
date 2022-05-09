

resource "aws_security_group" "public" {
  name        = "public-sg-${var.development_environment}"
  description = "Allow ALB TLS inbound traffic"
  vpc_id      = aws_vpc.head.id

  ingress {
    description      = "TLS from VPC"
    from_port        = 22
    to_port          = 22
    protocol         = "tcp"
    cidr_blocks      = [aws_vpc.head.cidr_block]

  }

  dynamic "ingress" {
    for_each = var.ingress_cidr_blocks
    content {
      description      = "http access"
      from_port        = 80
      to_port          = 80
      protocol         = "tcp"
      cidr_blocks      = [ingress.value]
    }

  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "test-public-sg"
    ManagedBy = "terraform"
    Environment = var.development_environment
  }

}
