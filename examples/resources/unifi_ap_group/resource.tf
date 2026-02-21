resource "unifi_ap_group" "all_aps" {
  name         = "All Access Points"
  for_wlanconf = true # Required for use with Wireless Networks (SSIDs)
}
