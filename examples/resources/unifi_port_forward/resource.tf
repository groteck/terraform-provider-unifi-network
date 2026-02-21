resource "unifi_port_forward" "web_server" {
  name     = "HTTPS Inbound"
  protocol = "tcp"
  src      = "any"
  dst_port = "443"
  fwd      = "192.168.1.50" # Internal IP
  fwd_port = "443"
  enabled  = true
}
