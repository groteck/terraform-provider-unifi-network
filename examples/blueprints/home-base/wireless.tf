resource "unifi_ap_group" "indoor" {
  name         = "Indoor APs"
  for_wlanconf = true
}

resource "unifi_wlan" "trusted" {
  name         = "HomeSecure"
  passphrase   = "trusted-secret"
  security     = "wpapsk"
  ap_group_ids = [unifi_ap_group.indoor.id]
  network_id   = unifi_network.trusted.id
}

resource "unifi_wlan" "iot" {
  name         = "HomeIOT"
  passphrase   = "iot-secret"
  security     = "wpapsk"
  ap_group_ids = [unifi_ap_group.indoor.id]
  network_id   = unifi_network.iot.id
}

resource "unifi_wlan" "guest" {
  name          = "HomeGuest"
  passphrase    = "welcome-guest"
  security      = "wpapsk"
  ap_group_ids  = [unifi_ap_group.indoor.id]
  network_id    = unifi_network.guest.id
  user_group_id = unifi_user_group.guest_limit.id
}


resource "unifi_wlan" "iot" {
  name            = "HomeIOT"
  passphrase      = "iot-secret"
  security        = "wpapsk"
  ap_group_ids    = [unifi_ap_group.indoor.id]
  network_conf_id = unifi_network.iot.id
}

resource "unifi_wlan" "guest" {
  name            = "HomeGuest"
  passphrase      = "welcome-guest"
  security        = "wpapsk"
  ap_group_ids    = [unifi_ap_group.indoor.id]
  network_conf_id = unifi_network.guest.id
  user_group_id   = unifi_user_group.guest_limit.id
}
