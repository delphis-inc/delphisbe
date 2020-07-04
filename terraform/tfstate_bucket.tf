resource "aws_s3_bucket" "tfstate_bucket" {
  bucket = "chatham-terraform"
  region = var.aws_region
}