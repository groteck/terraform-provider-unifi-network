package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRADIUSProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRADIUSProfileResourceConfig("Test RADIUS Profile"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_radius_profile.test", "name", "Test RADIUS Profile"),
					resource.TestCheckResourceAttr("unifi_radius_profile.test", "auth_servers.0.ip", "1.1.1.1"),
				),
			},
		},
	})
}

func testAccRADIUSProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "unifi_radius_profile" "test" {
  name = %[2]q
  auth_servers = [
    {
      ip     = "1.1.1.1"
      port   = 1812
      secret = "secret123"
    }
  ]
}
`, getProviderConfig(), name)
}
