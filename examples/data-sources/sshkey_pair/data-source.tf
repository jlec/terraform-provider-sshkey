terraform {
  required_version = "~>1.4"

  required_providers {
    turing-pi-bmc = {
      source  = "jlec.de/dev/sshkey"
      version = ">=0.1"
    }
  }
}

data "sshkey_pair" "example" {
  type = "rsa"
}

output "example" {
  value = data.sshkey_pair.example
}
