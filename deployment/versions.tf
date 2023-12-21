terraform {
  required_version = ">= 1.6.0"
  cloud {
    organization = "jlec-devops"

    workspaces {
      name = "terraform_provider_sshkey"
    }
  }
}
