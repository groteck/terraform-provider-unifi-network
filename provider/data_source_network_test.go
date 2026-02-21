package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkDataSourceConfig("Default"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.unifi_network.test", "id"),
					resource.TestCheckResourceAttr("data.unifi_network.test", "name", "Default"),
				),
			},
		},
	})
}

func testAccNetworkDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

data "unifi_network" "test" {
  name = %[2]q
}
`, getProviderConfig(), name)
}
