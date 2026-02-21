resource "unifi_user" "nas" {
  mac         = "00:11:22:33:44:55"
  name        = "Synology NAS"
  note        = "Main storage server"
  fixed_ip    = "192.168.1.10"
  use_fixedip = true
}
