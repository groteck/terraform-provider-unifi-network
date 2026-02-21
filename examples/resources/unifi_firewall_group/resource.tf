resource "unifi_firewall_group" "blacklist" {
  name          = "Malicious IPs"
  group_type    = "address-group"
  group_members = ["1.2.3.4", "5.6.7.8"]
}
