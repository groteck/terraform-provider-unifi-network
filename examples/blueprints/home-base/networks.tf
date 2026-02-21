resource "unifi_network" "trusted" {
  name    = "Trusted"
  vlan_id = 10
  subnet  = "192.168.10.1/24"
  purpose = "corporate"
}

resource "unifi_network" "iot" {
  name    = "IOT"
  vlan_id = 20
  subnet  = "192.168.20.1/24"
  purpose = "corporate"
}

resource "unifi_network" "guest" {
  name    = "Guest"
  vlan_id = 30
  subnet  = "192.168.30.1/24"
  purpose = "guest"
}
