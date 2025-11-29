package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), New("1.0.0"), providerserver.ServeOpts{
		Address: "registry.terraform.io/cloudflare/cloudflare-tunnel",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
