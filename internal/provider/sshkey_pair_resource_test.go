package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHKeyPairResource(t *testing.T) {
	t.Parallel()

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
