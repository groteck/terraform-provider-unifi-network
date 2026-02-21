package provider

import (
	"github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrafficRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrafficRuleResourceConfig("BLOCK", "INTERNET", "Test Traffic Rule"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_traffic_rule.test", "action", "BLOCK"),
					resource.TestCheckResourceAttr("unifi_traffic_rule.test", "matching_target", "INTERNET"),
					resource.TestCheckResourceAttr("unifi_traffic_rule.test", "description", "Test Traffic Rule"),
				),
			},
		},
	})
}

func testAccTrafficRuleResourceConfig(action, target, description string) string {
	return fmt.Sprintf(`
%s

resource "unifi_traffic_rule" "test" {
  action          = %[2]q
  matching_target = %[3]q
  description     = %[4]q
}
`, getProviderConfig(), action, target, description)
}
