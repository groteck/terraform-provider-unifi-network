package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccNetworkResourceConfig("Test Network", 100, "192.168.100.1/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "name", "Test Network"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "100"),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "192.168.100.1/24"),
				),
			},
			// Update and Read
			{
				Config: testAccNetworkResourceConfig("Updated Network", 101, "192.168.101.1/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_network.test", "name", "Updated Network"),
					resource.TestCheckResourceAttr("unifi_network.test", "vlan_id", "101"),
					resource.TestCheckResourceAttr("unifi_network.test", "subnet", "192.168.101.1/24"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPreCheck(t *testing.T) {
	// Add environment variable checks if necessary
}

func testAccNetworkResourceConfig(name string, vlan int, subnet string) string {
	return fmt.Sprintf(`
%s

resource "unifi_network" "test" {
  name    = %[2]q
  vlan_id = %[3]d
  subnet  = %[4]q
}
`, getProviderConfig(), name, vlan, subnet)
}
