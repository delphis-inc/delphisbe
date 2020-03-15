[
  {
    "name": "delphis-app",
    "image": "${app_image}",
    "cpu": ${fargate_cpu},
    "environment": [{
      "name": "DELPHIS_ENV",
      "value": "staging"
    }],
    "memory": ${fargate_memory},
    "networkMode": "awsvpc",
    "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/delphis-app",
          "awslogs-region": "${aws_region}",
          "awslogs-stream-prefix": "ecs"
        }
    },
    "portMappings": [
      {
        "containerPort": ${app_port},
        "hostPort": ${app_port}
      }
    ]
  }
]
