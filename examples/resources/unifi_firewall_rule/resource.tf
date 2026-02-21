resource "unifi_firewall_rule" "drop_iot_to_lan" {
  name     = "Drop IOT to LAN"
  ruleset  = "LAN_IN" # Controls traffic passing between networks
  action   = "drop"
  protocol = "all"
  enabled  = true

  # Matching connection states
  state_new         = true
  state_established = true
  state_invalid     = false
  state_related     = true

  # Rule order (lower numbers processed first)
  rule_index = 2000
}
