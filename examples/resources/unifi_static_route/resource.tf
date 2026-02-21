resource "unifi_static_route" "lab_gateway" {
  name     = "Lab Route"
  network  = "10.0.0.0/24"
  nexthop  = "192.168.1.254"
  distance = 1
  type     = "static-route"
}
