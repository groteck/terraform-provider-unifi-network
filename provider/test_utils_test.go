package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "unifi" {
  host           = "https://localhost:8443"
  username       = "admin"
  password       = "password123"
  allow_insecure = true
  is_standalone  = true
}
`
	providerConfigToken = `
provider "unifi" {
  host           = "https://localhost:8443"
  api_key        = "tf-test-token-12345"
  allow_insecure = true
  is_standalone  = true
}
`
)

func getProviderConfig() string {
	// You can add logic here to toggle between token and password
	// For now, we'll default to the one requested or both in different steps
	return providerConfigToken
}

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"unifi": providerserver.NewProtocol6WithError(New("test")()),
	}
)
