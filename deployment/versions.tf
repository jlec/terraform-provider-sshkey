terraform {
  required_version = ">= 1.9.0"
  cloud {
    organization = "jlec-devops"

    workspaces {
      name = "terraform_provider_sshkey"
    }
  }
}
