resource "unifi_radius_profile" "corporate" {
  name = "Corporate Auth"
  auth_servers = [
    {
      ip     = "192.168.1.10"
      port   = 1812
      secret = "radius-secret"
    }
  ]
}
