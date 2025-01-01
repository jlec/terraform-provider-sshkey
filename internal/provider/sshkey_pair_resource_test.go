/*
Copyright 2022-2025 Justin Lecher

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
