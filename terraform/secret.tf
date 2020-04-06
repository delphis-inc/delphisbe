# resource "aws_secretsmanager_secret" "chatham-staging-aurora-pgsql-pass" {
#     name = "chatham-staging-auorora-db-pass"
# }

# resource "aws_secretsmanager_secret_version" "db-pass-val" {
#     secret_id     = "${aws_secretsmanager_secret.chatham-staging-aurora-pgsql-pass.id}"
#     secret_string = "${random_string.db_master_pass.result}"
# }