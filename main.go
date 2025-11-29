package main

import (
	"context"
	"log"

	"terraform-provider-cloudflare-tunnel/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New("1.0.0"), providerserver.ServeOpts{
		Address: "registry.terraform.io/cloudflare/cloudflare-tunnel",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
