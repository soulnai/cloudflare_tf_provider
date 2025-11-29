package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudflareTunnel_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "cloudflare-tunnel" "test" {
  name       = "tf-provider-test-tunnel"
  tunnel_token = "AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg=" # 32 bytes base64
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare-tunnel.test", "name", "tf-provider-test-tunnel"),
					resource.TestCheckResourceAttrSet("cloudflare-tunnel.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cloudflare-tunnel.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Secret is sensitive and not returned by API, so we skip verification for it
				ImportStateVerifyIgnore: []string{"secret", "tunnel_token"},
			},
		},
	})
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cloudflare": providerserver.NewProtocol6WithError(New("test")()),
}

const providerConfig = `
provider "cloudflare" {
  api_token  = "test-token"  # Usually picked up from env var CLOUDFLARE_API_TOKEN
  account_id = "test-account"
  base_url   = "http://localhost:8080" # Mock server or real URL
}
`
