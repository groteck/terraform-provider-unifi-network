resource "unifi_static_dns" "home_assistant" {
  key         = "ha.local"
  value       = "192.168.1.20"
  record_type = "A"
  enabled     = true
}
