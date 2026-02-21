resource "unifi_network" "homelab" {
  name    = "Homelab"
  vlan_id = 40
  subnet  = "192.168.40.1/24"
  purpose = "corporate"
}

# Allow Trusted to access Homelab
resource "unifi_firewall_rule" "trusted_to_lab" {
  name     = "Allow Trusted to Lab"
  ruleset  = "LAN_IN"
  action   = "accept"
  protocol = "all"
  enabled  = true

  src_network_id   = unifi_network.trusted.id # Defined in home-base
  src_network_type = "NETv4"
  dst_network_id   = unifi_network.homelab.id
  dst_network_type = "NETv4"
}

# Block Homelab from accessing Trusted
resource "unifi_firewall_rule" "block_lab_to_trusted" {
  name     = "Block Lab to Trusted"
  ruleset  = "LAN_IN"
  action   = "drop"
  protocol = "all"
  enabled  = true

  src_network_id   = unifi_network.homelab.id
  src_network_type = "NETv4"
  dst_network_id   = unifi_network.trusted.id
  dst_network_type = "NETv4"
}

# Port forward for HomeLab service
resource "unifi_port_forward" "lab_server" {
  name     = "Lab HTTP"
  protocol = "tcp"
  src      = "any"
  dst_port = "80"
  fwd      = "192.168.40.10"
  fwd_port = "80"
}

# Static IP for Lab Server
resource "unifi_user" "lab_server" {
  mac         = "aa:bb:cc:dd:ee:ff"
  name        = "Lab Server"
  fixed_ip    = "192.168.40.10"
  use_fixedip = true
  network_id  = unifi_network.homelab.id
}
