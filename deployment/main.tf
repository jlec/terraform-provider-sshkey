locals {
}

data "tfe_organization" "jlec" {
  name = "jlec"
}

resource "tfe_registry_provider" "sshkey" {
  organization = data.tfe_organization.jlec.name

  registry_name = "public"
  namespace     = "jlec"
  name          = "sshkey"
}
