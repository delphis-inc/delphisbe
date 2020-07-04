# NOTE: Leaving this here because it can be helpful, however
# two reasons to not create a DB with terraform:
# 1. It stores the password in plain text in the terraform.tfstate
#    file which we are checking in. I don't want to have this in
#    an open source project.
# 2. As we saw with dynamo tables, it's quite easy to eff something
#    up compeltely accidentally and destroy the entire DB. Given
#    that creation of new databases is pretty rare I'll check in
#    a script instead. Then we can take the ARN and put that into
#    the IAM pieces of Terraform instead.

# resource "random_string" "db_master_pass" {
#     length           = 40
#     override_special = "!#$%^&*()`_=+[]{}<>:?"
#     special          = true
#     min_special      = 5
#     keepers          = {
#         pass_version = 1
#     }
# }

# resource "aws_rds_cluster" "postgresql" {
#     cluster_identifier      = "chatham-staging-aurora-pgsql"
#     engine                  = "aurora-postgresql"
#     # Setting this to single AZ for staging service.
#     availability_zones      = ["us-west-2a"]
#     database_name           = "chatham-staging"
#     master_username         = "postgres"
#     # I know this will still be written to tfstate in raw text. I don't
#     # know how to get around it though... Perhaps we create this by hand
#     # instead? Probably a TODO for down the road.
#     password                = "${random_string.db_master_pass.result}"
#     backup_retention_period = 5
#     preferred_backup_window = "07:00-09:00"
# }

resource "aws_rds_cluster_instance" "cluster_instances-staging" {
  count      = 1
  identifier = "chatham-staging-aurora-psgql-${count.index}"
  # Copy this manually from the account if you want to change!
  cluster_identifier = "chatham-staging-aurora-pgsql"
  # Smallest we can create...
  instance_class = "db.r4.large"
  engine         = "aurora-postgresql"
}