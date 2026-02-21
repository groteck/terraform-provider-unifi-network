resource "unifi_user_group" "guest_limit" {
  name           = "Guest Bandwidth"
  download_limit = 5000
  upload_limit   = 1000
}
