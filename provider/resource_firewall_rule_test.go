package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleResourceConfig("Test Rule", "LAN_LOCAL", "accept", "all"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "name", "Test Rule"),
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "ruleset", "LAN_LOCAL"),
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "action", "accept"),
					resource.TestCheckResourceAttrSet("unifi_firewall_rule.test", "rule_index"),
				),
			},
			{
				Config: testAccFirewallRuleResourceConfig("Updated Rule", "LAN_LOCAL", "drop", "tcp"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "name", "Updated Rule"),
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "action", "drop"),
					resource.TestCheckResourceAttr("unifi_firewall_rule.test", "protocol", "tcp"),
				),
			},
		},
	})
}

func testAccFirewallRuleResourceConfig(name, ruleset, action, protocol string) string {
	return fmt.Sprintf(`
%s

resource "unifi_firewall_rule" "test" {
  name     = %[2]q
  ruleset  = %[3]q
  action   = %[4]q
  protocol = %[5]q
  enabled  = true
  
  src_network_type = "ADDRv4"
  dst_network_type = "ADDRv4"

  state_new         = true
  state_established = true
  state_invalid     = false
  state_related     = true
  ipsec             = ""
  rule_index        = 2000
  logging           = false
}

`, getProviderConfig(), name, ruleset, action, protocol)
}
