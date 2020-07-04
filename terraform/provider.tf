# provider.tf

# Specify the provider and access details
provider "aws" {
  shared_credentials_file = "$HOME/.aws/credentials"
  profile                 = "delphis"
  region                  = var.aws_region
}

terraform {
  backend "s3" {
    encrypt = true
    region  = "us-west-2"
  }
}