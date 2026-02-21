resource "unifi_port_profile" "iot_port" {
  name              = "IOT Port"
  native_network_id = unifi_network.iot.id
  forward           = "all"
}
