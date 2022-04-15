resource "aws_lb_target_group" "main" {
  name        = var.lb_target_group_names[0]
  target_type = var.lb_target_type # "lambda"
}

resource "aws_lambda_permission" "with_lb" {
  statement_id  = "AllowExecutionFromlbtoLambda"
  action        = "lambda:InvokeFunction"
  function_name = var.target_ids_arn[0]
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.main.arn 
  # qualifier     = aws_lambda_alias.lambda_alias.name
}


resource "aws_lb_target_group_attachment" "main" {
  target_group_arn = aws_lb_target_group.main.arn
  target_id        = var.target_ids_arn[0] # example: aws_lambda_function.main.arn
  depends_on       = [aws_lambda_permission.with_lb]
}


# # *****   START = second lb group ***** 

# resource "aws_lb_target_group" "sqs_group" {
#   name        = var.lb_target_group_names[1]
#   target_type = var.lb_target_type # "lambda"
# }

# resource "aws_lambda_permission" "with_lb_sqs" {
#   statement_id  = "AllowExecutionFromlbtoLambda"
#   action        = "lambda:InvokeFunction"
#   function_name = var.target_ids_arn[1]
#   principal     = "elasticloadbalancing.amazonaws.com"
#   source_arn    = aws_lb_target_group.sqs_group.arn 
#   # qualifier     = aws_lambda_alias.lambda_alias.name
# }


# resource "aws_lb_target_group_attachment" "sqs_group" {
#   target_group_arn = aws_lb_target_group.sqs_group.arn
#   target_id        = var.target_ids_arn[1] # example: aws_lambda_function.main.arn
#   depends_on       = [aws_lambda_permission.with_lb_sqs]
# }

# # *****   END = second lb group ***** 


resource "aws_lb" "main" {
  name               = var.load_balancer_name
  internal           = false
  load_balancer_type = "application"
  security_groups    = var.security_groups #[module.my_vpc.security_group]
  subnets            = var.subnet_ids  #module.my_vpc.subnet_ids

  enable_deletion_protection = false

  tags = {
    Environment = var.development_environment
  }

  
}

resource "aws_lb_listener" "main" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "forward"
    target_group_arn = aws_lb_target_group.main.arn
  }
}

# resource "aws_lb_listener_rule" "firehose_lambda_target" {
#   listener_arn = aws_lb_listener.main.arn

#   action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.sqs_group.arn
#   }

#   condition {
#     query_string {
#       key   = "flag"
#       value = "Y"
#     }
#   }
# }

# resource "aws_lb_listener_rule" "sqs_lambda_target" {
#   listener_arn = aws_lb_listener.main.arn

#   action {
#     type             = "forward"
#     target_group_arn = aws_lb_target_group.main.arn
#   }

#   condition {
#     query_string {
#       key   = "flag"
#       value = "A"
#     }
#   }
# }