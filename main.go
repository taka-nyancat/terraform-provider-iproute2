package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/taka-nyancat/terraform-provider-iproute2/internal/provider"
)


func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/taka-nyancat/iproute2",
	})
	if err != nil {
		log.Fatal(err)
	}
}
