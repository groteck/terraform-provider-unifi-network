package provider

import (
	"github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStaticRouteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteResourceConfig("Test Route", "10.0.0.0/24", "192.168.1.254", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_route.test", "name", "Test Route"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "network", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "distance", "1"),
				),
			},
		},
	})
}

func testAccStaticRouteResourceConfig(name, network, nexthop string, distance int) string {
	return fmt.Sprintf(`
%s

resource "unifi_static_route" "test" {
  name     = %[2]q
  network  = %[3]q
  nexthop  = %[4]q
  distance = %[5]d
  type     = "static-route"
}
`, getProviderConfig(), name, network, nexthop, distance)
}
