resource "unifi_user_group" "guest_limit" {
  name           = "Guest Limits"
  download_limit = 5000 # Kbps (5 Mbps)
  upload_limit   = 1000 # Kbps (1 Mbps)
}
