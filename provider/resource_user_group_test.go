package provider

import (
	"github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupResourceConfig("Test Group", 1000, 1000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user_group.test", "name", "Test Group"),
					resource.TestCheckResourceAttr("unifi_user_group.test", "download_limit", "1000"),
					resource.TestCheckResourceAttr("unifi_user_group.test", "upload_limit", "1000"),
				),
			},
		},
	})
}

func testAccUserGroupResourceConfig(name string, dl, ul int) string {
	return fmt.Sprintf(`
%s

resource "unifi_user_group" "test" {
  name           = %[2]q
  download_limit = %[3]d
  upload_limit   = %[4]d
}
`, getProviderConfig(), name, dl, ul)
}
