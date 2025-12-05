terraform {
  required_version = ">= 1.10.0"
  cloud {}
  required_providers {
    tfe = {
      source  = "hashicorp/tfe"
      version = ">=0.71.0, <1.0.0"
    }
  }
}

provider "tfe" {
  # Configuration options
}
