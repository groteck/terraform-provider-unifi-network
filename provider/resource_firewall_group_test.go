package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallGroupResourceConfig("Test Group", "address-group", []string{"192.168.1.1", "192.168.1.2"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_firewall_group.test", "name", "Test Group"),
					resource.TestCheckResourceAttr("unifi_firewall_group.test", "group_type", "address-group"),
					resource.TestCheckResourceAttr("unifi_firewall_group.test", "group_members.#", "2"),
				),
			},
		},
	})
}

func testAccFirewallGroupResourceConfig(name, groupType string, members []string) string {
	membersStr := ""
	for i, m := range members {
		membersStr += fmt.Sprintf("%q", m)
		if i < len(members)-1 {
			membersStr += ", "
		}
	}
	return fmt.Sprintf(`
%s

resource "unifi_firewall_group" "test" {
  name          = %[2]q
  group_type    = %[3]q
  group_members = [%[4]s]
}
`, getProviderConfig(), name, groupType, membersStr)
}
