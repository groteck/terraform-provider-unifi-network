resource "unifi_network" "iot" {
  name    = "IOT"
  vlan_id = 20
  subnet  = "192.168.20.1/24"
  purpose = "corporate" # Usually 'corporate' for internal VLANs
}
