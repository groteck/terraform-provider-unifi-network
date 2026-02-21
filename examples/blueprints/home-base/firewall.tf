resource "unifi_firewall_rule" "drop_iot_to_trusted" {
  name     = "Drop IOT to Trusted"
  ruleset  = "LAN_IN"
  action   = "drop"
  protocol = "all"
  enabled  = true

  src_network_id   = unifi_network.iot.id
  src_network_type = "NETv4"
  dst_network_id   = unifi_network.trusted.id
  dst_network_type = "NETv4"

  state_new         = true
  state_established = false
  state_invalid     = true
  state_related     = false
}

resource "unifi_firewall_rule" "drop_guest_to_trusted" {
  name     = "Drop Guest to Trusted"
  ruleset  = "LAN_IN"
  action   = "drop"
  protocol = "all"
  enabled  = true

  src_network_id   = unifi_network.guest.id
  src_network_type = "NETv4"
  dst_network_id   = unifi_network.trusted.id
  dst_network_type = "NETv4"
}
