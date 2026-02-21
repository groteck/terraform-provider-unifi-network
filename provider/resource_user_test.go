package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig("00:11:22:33:44:55", "Test Device", "192.168.1.100"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "mac", "00:11:22:33:44:55"),
					resource.TestCheckResourceAttr("unifi_user.test", "name", "Test Device"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", "192.168.1.100"),
					resource.TestCheckResourceAttr("unifi_user.test", "use_fixedip", "true"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(mac, name, ip string) string {
	return fmt.Sprintf(`
%s

resource "unifi_user" "test" {
  mac         = %[2]q
  name        = %[3]q
  fixed_ip    = %[4]q
  use_fixedip = true
}
`, getProviderConfig(), mac, name, ip)
}
