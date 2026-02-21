package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWLANResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWLANResourceConfig("TestSSID", "password123"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_wlan.test", "name", "TestSSID"),
					resource.TestCheckResourceAttr("unifi_wlan.test", "security", "wpapsk"),
				),
			},
		},
	})
}

func testAccWLANResourceConfig(name, passphrase string) string {
	return fmt.Sprintf(`
%s

resource "unifi_ap_group" "test_wlan" {
  name         = "WLAN Test Group"
  for_wlanconf = false
}

resource "unifi_wlan" "test" {
  name         = %[2]q
  passphrase   = %[3]q
  security     = "wpapsk"
  ap_group_ids = [unifi_ap_group.test_wlan.id]
}
`, getProviderConfig(), name, passphrase)
}
