
# this module will create application load balancer
module "my_ec2" {
  source = "../modules/ec2"

  lb_target_group_names = ["lambda-target-group-firehose-sqs"]
  lb_target_type = "lambda"
  target_ids_arn = [aws_lambda_function.main.arn]
  load_balancer_name = var.load_balancer_name
  security_groups = [module.my_vpc.security_group]
  subnet_ids = module.my_vpc.subnet_ids
  development_environment = var.development_environment

}


