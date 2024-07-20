data "aws_iam_policy_document" "ecs_node" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_node" {
  assume_role_policy = data.aws_iam_policy_document.ecs_node
}

resource "aws_iam_role_policy_attachment" "ecs_node_role_policy" {
  role       = aws_iam_role.ecs_node.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs_node" {
  path = "/ecs/instance/"
  role = aws_iam_role.ecs_node.name
}

# --- ECS Node SG ---

resource "aws_security_group" "ecs_node_sg" {
  vpc_id = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_kms_key" "garrettdavis_dev" {
  description             = "garrettdavis_dev"
  deletion_window_in_days = 7
}

resource "aws_cloudwatch_log_group" "garrettdavis_dev" {
  name = "garrettdavis_dev"
}

resource "aws_ecs_cluster" "garrettdavis_dev" {
  name = "garrettdavis_dev"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  configuration {
    execute_command_configuration {
      kms_key_id = aws_kms_key.garrettdavis_dev.arn
      logging    = "OVERRIDE"

      log_configuration {
        cloud_watch_encryption_enabled = true
        cloud_watch_log_group_name     = aws_cloudwatch_log_group.garrettdavis_dev.name
      }
    }
  }
}

resource "aws_autoscaling_group" "garrettdavis_dev" {
  name     = "garrettdavis_dev"
  max_size = 3
  min_size = 1
}

resource "aws_ecs_capacity_provider" "garrettdavis_dev" {
  name = "garrettdavis_dev"

  auto_scaling_group_provider {
    auto_scaling_group_arn         = aws_autoscaling_group.garrettdavis_dev.arn
    managed_termination_protection = "ENABLED"

    managed_scaling {
      maximum_scaling_step_size = 1000
      minimum_scaling_step_size = 1
      status                    = "ENABLED"
      target_capacity           = 10
    }
  }
}

resource "aws_ecs_cluster_capacity_providers" "garrettdavis_dev" {
  cluster_name = aws_ecs_cluster.garrettdavis_dev.name

  capacity_providers = [
  ]
}
