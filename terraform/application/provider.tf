# defined aws provider
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.region
#   access_key = "provide_your_access_key" 
#   secret_key = "provide_your_secret_key"
}
