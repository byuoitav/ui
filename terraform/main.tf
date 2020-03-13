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

// pull all env vars out of ssm
data "aws_ssm_parameter" "dev_couch_address" {
  name = "/env/dev-couch-address"
}

data "aws_ssm_parameter" "dev_couch_username" {
  name = "/env/dev-couch-username"
}

data "aws_ssm_parameter" "dev_couch_password" {
  name = "/env/dev-couch-password"
}

data "aws_ssm_parameter" "dev_hub_address" {
  name = "/env/dev-hub-address"
}

data "aws_ssm_parameter" "dev_code_service_address" {
  name = "/env/dev-code-service-address"
}

module "deployment" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "av-control-ui-dev"
  image          = "docker.pkg.github.com/byuoitav/ui/ui-dev"
  image_version  = "01e1354"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/ui"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["roomcontrol-dev.av.byu.edu", "rooms-dev.av.byu.edu"]
  container_env = {
    "DB_ADDRESS"       = data.aws_ssm_parameter.dev_couch_address.value
    "DB_USERNAME"      = data.aws_ssm_parameter.dev_couch_username.value
    "DB_PASSWORD"      = data.aws_ssm_parameter.dev_couch_password.value
    "STOP_REPLICATION" = "true"
    "CODE_SERVICE_URL" = data.aws_ssm_parameter.dev_code_service_address.value
    "HUB_ADDRESS"      = data.aws_ssm_parameter.dev_hub_address.value
  }
  container_args = [
    "--port", "8080",
    "--log-level", "2", // set log level to info
    // "--av-api", "av-api-prd.default.svc.cluster.local",
    "--av-api", "itb-1006-cp1.byu.edu:8000",
    "--lazarette", "lazarette-dev.default.svc.cluster.local",
    "--code-service", data.aws_ssm_parameter.dev_code_service_address.value,
  ]
}

// TODO prod
