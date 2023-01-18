terraform {
  required_version = ">= 1.3.0"
  cloud {
    organization = "jlec"

    workspaces {
      name = "terraform_provider_sshkey"
    }
  }
}
