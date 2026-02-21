package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAPGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAPGroupResourceConfig("Test AP Group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_ap_group.test", "name", "Test AP Group"),
				),
			},
		},
	})
}

func testAccAPGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "unifi_ap_group" "test" {
  name         = %[2]q
  for_wlanconf = false
}
`, getProviderConfig(), name)
}
