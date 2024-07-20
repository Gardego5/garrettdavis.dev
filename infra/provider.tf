terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "tf-state-20230722071359242500000001"
    key            = "state/garrettdavis_dev_2"
    region         = "us-west-2"
    encrypt        = true
    kms_key_id     = "alias/terraform-state"
    dynamodb_table = "tf-state-20230722071359242500000001"
  }
}

provider "aws" {
  region = "us-west-2"
}
