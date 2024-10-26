terraform {
  required_version = ">= 1.9.0"

  required_providers {
    turing-pi-bmc = {
      source  = "jlec.de/dev/sshkey"
      version = ">=0.1"
    }
  }
}

resource "sshkey_pair" "example" {
  type = "rsa"
}
