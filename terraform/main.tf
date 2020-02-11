terraform {
  backend "s3" {
    bucket     = "terraform-state-storage-586877430255"
    lock_table = "terraform-state-lock-586877430255"
    region     = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "av-control-ui.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
}

module "deployment" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "av-control-ui-dev"
  image          = "docker.pkg.github.com/byuoitav/ui/amd64"
  image_version  = "v0.3.0"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/ui"

  // optional
  image_pull_secret = "github-docker-registry"
  public_url        = "rooms.dev.av.byu.edu"
  container_env = {
    "DB_USERNAME"      = var.db_username
    "DB_PASSWORD"      = var.db_password
    "DB_ADDRESS"       = var.db_address
    "CODE_SERVICE_URL" = var.code_service_url
    "HUB_ADDRESS"      = var.hub_address
  }
}

// TODO add the other route53 entry
