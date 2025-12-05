terraform {
  required_version = ">= 1.9.0"

  required_providers {
    sshkey = {
      source  = "jlec.de/dev/sshkey"
      version = ">=0.1"
    }
  }
}
provider "sshkey" {
  # example configuration here
}
