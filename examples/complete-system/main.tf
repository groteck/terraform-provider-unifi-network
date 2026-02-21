terraform {
  required_providers {
    unifi = {
      source = "jlopez/unifi"
    }
  }
}

provider "unifi" {
  # host = "https://192.168.1.1"
  # api_key = "..."
  # is_standalone = true
}

# --- Networks ---

resource "unifi_network" "trusted" {
  name    = "Trusted"
  vlan_id = 10
  subnet  = "192.168.10.1/24"
}

resource "unifi_network" "iot" {
  name    = "IOT"
  vlan_id = 20
  subnet  = "192.168.20.1/24"
}

resource "unifi_network" "guest" {
  name    = "Guest"
  vlan_id = 30
  subnet  = "192.168.30.1/24"
  purpose = "guest"
}

# --- Firewall Groups ---

resource "unifi_firewall_group" "trusted_networks" {
  name          = "Trusted Networks"
  group_type    = "address-group"
  group_members = ["192.168.10.0/24"]
}

# --- Firewall Rules ---

# Default Drop from IOT to Trusted
resource "unifi_firewall_rule" "block_iot_to_trusted" {
  name     = "Isolate IOT"
  ruleset  = "LAN_IN"
  action   = "drop"
  protocol = "all"

  src_network_id   = unifi_network.iot.id
  src_network_type = "NETv4"
  dst_network_id   = unifi_network.trusted.id
  dst_network_type = "NETv4"

  state_new         = true
  state_established = false
  state_invalid     = true
  state_related     = false
}

# Allow Established/Related back to IOT
resource "unifi_firewall_rule" "allow_trusted_to_iot" {
  name     = "Allow Trusted to IOT"
  ruleset  = "LAN_IN"
  action   = "accept"
  protocol = "all"

  src_network_id   = unifi_network.trusted.id
  src_network_type = "NETv4"
  dst_network_id   = unifi_network.iot.id
  dst_network_type = "NETv4"
}


# --- Wireless ---

resource "unifi_ap_group" "default" {
  name         = "Default AP Group"
  for_wlanconf = true
}

resource "unifi_wlan" "main" {
  name            = "SecureNet"
  passphrase      = "v3ry-s3cur3"
  security        = "wpapsk"
  ap_group_ids    = [unifi_ap_group.default.id]
  network_conf_id = unifi_network.trusted.id
}

resource "unifi_wlan" "iot" {
  name            = "IOTNet"
  passphrase      = "iot-only"
  security        = "wpapsk"
  ap_group_ids    = [unifi_ap_group.default.id]
  network_conf_id = unifi_network.iot.id
}

# --- Users & DHCP ---

resource "unifi_user" "home_server" {
  mac         = "11:22:33:44:55:66"
  name        = "Home Server"
  fixed_ip    = "192.168.10.10"
  use_fixedip = true
  network_id  = unifi_network.trusted.id
}
