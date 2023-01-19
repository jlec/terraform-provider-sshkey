package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSSHKeyPairResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSSHKeyPairResourceConfig("rsa"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sshkey_pair.test", "type", "rsa"),
					resource.TestCheckResourceAttr("sshkey_pair.test", "bits", "4096"),
				),
			},
			// ImportState testing
			// {
			// 	ResourceName:      "sshkey_pair.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// This is not normally necessary, but is here because this
			// 	// example code does not have an actual upstream service.
			// 	// Once the Read method is able to refresh information from
			// 	// the upstream service, this can be removed.
			// 	ImportStateVerifyIgnore: []string{"configurable_attribute"},
			// },
			// Update and Read testing
			// {
			// 	Config: testAccSSHKeyPairResourceConfig("ed25519"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("sshkey_pair.test", "comment", "fuzzy"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSSHKeyPairResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "sshkey_pair" "test" {
  type = %[1]q
}
`, configurableAttribute)
}
