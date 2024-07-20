resource "aws_ecr_repository" "garrettdavis_dev" {
  name                 = "garrettdavis_dev"
  image_tag_mutibility = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecs_cluster" "ecs_cluster" {
  name = "garrettdavis_dev"
}

resource "aws_cloudwatch_log_group" "garrettdavis_dev" {
  name              = "/ecs/garrettdavis_dev"
  retention_in_days = 14
}

resource "aws_ecs_task_definition" "garrettdavis_dev" {
  family = "garrettdavis_dev"
  container_definitions = jsonencode([
    {
      name        = "app"
      image       = "${aws_ecr_repository.garrettdavis_dev.repository_url}:latest"
      entryPoint  = []
      essential   = true
      networkMode = "awsvpc"
      portMappings = [{
        containerPort = 3000
        hostPort      = 80
      }]
      healthCheck = {
        command     = ["CMD-SHELL", "curl -f https://localhost:3000/ || exit 1"]
        interval    = 2, timeout = 5
        startPeriod = 2, retries = 5
      }
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-region        = "us-west-2"
          awslogs-group         = aws_cloudwatch_log_group.garrettdavis_dev.name
          awslogs-stream-prefix = "app"
        }
      }
    }
  ])
}

resource "aws_security_group" "ecs_task" {
  name_prefix = "ecs-task-sg-"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_service" "garrettdavis_dev" {
  name            = "app"
  cluster         = aws_ecs_cluster.garrettdavis_dev.id
  task_definition = aws_ecs_task_definition.garrettdavis_dev.arn
  desired_count   = 1

  network_configuration {
    security_groups = [aws_security_group.ecs_task.id]
    subnets         = aws_subnet.public[*].id
  }

  capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_providier.garrettdavis_dev.name
    base              = 1
    weight            = 100
  }

  ordered_placement_strategy {
    type  = "spread"
    field = "attribute:ecs.availability-zone"
  }

  lifecycle {
    ignore_changes = [desired_count]
  }
}
