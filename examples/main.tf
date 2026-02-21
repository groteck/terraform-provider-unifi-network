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

resource "unifi_network" "home" {
  name    = "Home"
  vlan_id = 10
  subnet  = "192.168.10.1/24"
}

resource "unifi_user_group" "unlimited" {
  name = "Unlimited"
}

resource "unifi_ap_group" "main" {
  name         = "Main APs"
  for_wlanconf = false
}

resource "unifi_wlan" "main" {
  name         = "MyWiFi"
  passphrase   = "secret123"
  ap_group_ids = [unifi_ap_group.main.id]
}

resource "unifi_firewall_group" "iot_ips" {
  name          = "IOT IPs"
  group_type    = "address-group"
  group_members = ["192.168.10.50", "192.168.10.51"]
}

resource "unifi_firewall_rule" "drop_iot" {
  name     = "Drop IOT"
  ruleset  = "LAN_IN"
  action   = "drop"
  protocol = "all"

  state_new         = true
  state_established = true
  ipsec             = ""
}

resource "unifi_user" "my_pc" {
  mac         = "00:11:22:33:44:55"
  name        = "My PC"
  fixed_ip    = "192.168.10.100"
  use_fixedip = true
}
