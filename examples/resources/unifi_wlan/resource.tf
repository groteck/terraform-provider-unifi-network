resource "unifi_wlan" "home_wifi" {
  name         = "HomeWiFi"
  passphrase   = "securepassword"
  security     = "wpapsk"
  ap_group_ids = [unifi_ap_group.all_aps.id]
}
