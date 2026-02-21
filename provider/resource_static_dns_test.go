package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStaticDNSResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticDNSResourceConfig("home.local", "192.168.1.10"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_dns.test", "key", "home.local"),
					resource.TestCheckResourceAttr("unifi_static_dns.test", "value", "192.168.1.10"),
				),
			},
		},
	})
}

func testAccStaticDNSResourceConfig(key, value string) string {
	return fmt.Sprintf(`
%s

resource "unifi_static_dns" "test" {
  key   = %[2]q
  value = %[3]q
  ttl   = 3600
}
`, getProviderConfig(), key, value)
}
