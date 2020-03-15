# ECS task execution role data
data "aws_iam_policy_document" "ecs_task_execution_role" {
  version = "2012-10-17"
  statement {
    sid = ""
    effect = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

# ECS task execution role
resource "aws_iam_role" "ecs_task_execution_role" {
  name               = var.ecs_task_execution_role_name
  assume_role_policy = data.aws_iam_policy_document.ecs_task_execution_role.json
}

# ECS task execution role policy attachment
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "ecs_task_role" {
  name = "delphis-ecsTaskRole"
 
  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "ecs-tasks.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
}
 
resource "aws_iam_policy" "dynamodb" {
  name        = "delphis-task-policy-dynamodb"
  description = "Policy that allows access to DynamoDB"
 
 policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
       {
           "Effect": "Allow",
           "Action": [
               "dynamodb:UpdateTimeToLive",
               "dynamodb:PutItem",
               "dynamodb:BatchWriteItem",
               "dynamodb:DescribeTable",
               "dynamodb:ListTables",
               "dynamodb:DeleteItem",
               "dynamodb:GetItem",
               "dynamodb:BatchGetItem",
               "dynamodb:Scan",
               "dynamodb:Query",
               "dynamodb:UpdateItem"
           ],
           "Resource": "*"
       }
   ]
}
EOF
}
 
resource "aws_iam_role_policy_attachment" "ecs-task-role-policy-attachment" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.dynamodb.arn
}

// TODO: This should limit to just the speicfic secret that we use for delphis.
resource "aws_iam_policy" "secrets-manager" {
  name        = "delphis-task-policy-secrets-manager"
  description = "Policy that allows access to Secrets-manager"
 
 policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
       {
           "Effect": "Allow",
           "Action": [
              "secretsmanager:Describe*",
              "secretsmanager:Get*",
              "secretsmanager:List*" 
           ],
           "Resource": "*"
       }
   ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-task-role-policy-attachment-secrets-manager" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.secrets-manager.arn
}