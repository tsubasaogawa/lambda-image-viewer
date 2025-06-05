terraform {
  required_version = "~> 1.12.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.89"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"

  default_tags {
    tags = {
      repository = "tsubasaogawa/lambda-image-viewer"
    }
  }
}

provider "aws" {
  region = "us-east-1"
  alias  = "n_virginia"

  default_tags {
    tags = {
      repository = "tsubasaogawa/lambda-image-viewer"
    }
  }
}
