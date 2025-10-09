terraform {
  required_version = ">= 1.13.0"
  cloud {
    organization = "jlec-devops"

    workspaces {
      name = "terraform_provider_sshkey"
    }
  }
}
