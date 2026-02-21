package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAPGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAPGroupDataSourceConfig("Test AP Group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.unifi_ap_group.test", "id"),
					resource.TestCheckResourceAttr("data.unifi_ap_group.test", "name", "Test AP Group"),
				),
			},
		},
	})
}

func testAccAPGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "unifi_ap_group" "test" {
  name         = %[2]q
  for_wlanconf = false
}

data "unifi_ap_group" "test" {
  name = unifi_ap_group.test.name
}
`, getProviderConfig(), name)
}
