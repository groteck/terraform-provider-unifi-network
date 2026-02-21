resource "unifi_traffic_rule" "block_tiktok" {
  description     = "Block Social Media"
  action          = "BLOCK"
  matching_target = "APP" # Can also be DOMAIN, IP, or INTERNET
  enabled         = true
}
