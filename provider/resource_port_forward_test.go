package provider

import (
	"github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPortForwardResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardResourceConfig("Test Forward", "tcp", "80", "192.168.1.100", "8080"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_forward.test", "name", "Test Forward"),
					resource.TestCheckResourceAttr("unifi_port_forward.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("unifi_port_forward.test", "fwd", "192.168.1.100"),
				),
			},
		},
	})
}

func testAccPortForwardResourceConfig(name, protocol, dstPort, fwd, fwdPort string) string {
	return fmt.Sprintf(`
%s

resource "unifi_port_forward" "test" {
  name     = %[2]q
  protocol = %[3]q
  dst_port = %[4]q
  fwd      = %[5]q
  fwd_port = %[6]q
}
`, getProviderConfig(), name, protocol, dstPort, fwd, fwdPort)
}
