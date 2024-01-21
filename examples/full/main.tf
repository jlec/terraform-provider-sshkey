terraform {
  required_version = ">= 1.3.0"
  required_providers {
    sshkey = {
      # source = "app.terraform.io/jlec/sshkey"
      source  = "jlec.de/dev/sshkey"
      version = ">= 0.0.1"
    }
    local = {
      source  = "hashicorp/local"
      version = "~>2.0"
    }
  }
}

resource "sshkey_pair" "rsa" {
  type = "rsa"
  bits = 2048
}

resource "sshkey_pair" "ed25519" {
  type    = "ed25519"
  comment = "fuzzy"
}

resource "sshkey_pair" "ecdsa" {
  type    = "ecdsa"
  comment = "admin@example.com"
}

resource "local_file" "rsa" {
  filename        = "id_rsa"
  content         = sshkey_pair.rsa.private_key
  file_permission = 0600
}

resource "local_file" "rsa_pub" {
  filename        = "id_rsa.pub"
  content         = sshkey_pair.rsa.public_key
  file_permission = 0600
}

resource "local_file" "ed25519" {
  filename        = "id_ed25519"
  content         = sshkey_pair.ed25519.private_key
  file_permission = 0600
}

resource "local_file" "ed25519_pub" {
  filename        = "id_ed25519.pub"
  content         = sshkey_pair.ed25519.public_key
  file_permission = 0600
}

# output "example_fingerprint_md5" {
#   value = sshkey.example[*].fingerprint_md5
# }

# output "big_public_key" {
#   value = sshkey.big.public_key
# }

# output "small_fingerprint_sha256" {
#   value = sshkey.small.fingerprint_sha256
# }

# output "small_public_key" {
#   value = sshkey.small.public_key
# }
