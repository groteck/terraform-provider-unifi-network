package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPortProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortProfileResourceConfig("Test Profile"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_profile.test", "name", "Test Profile"),
				),
			},
		},
	})
}

func testAccPortProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
%s

data "unifi_network" "default" {
  name = "Default"
}

resource "unifi_port_profile" "test" {
  name              = %[2]q
  native_network_id = data.unifi_network.default.id
  forward           = "all"
}
`, getProviderConfig(), name)
}
