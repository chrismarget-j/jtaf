// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated"
	jtafProvider "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const registryAddr = "%s/juniper/%s"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()
	opts := providerserver.ServeOpts{
		Address: fmt.Sprintf(registryAddr, generated.RegistryHost, generated.ProviderName),
		Debug:   debug,
	}

	err := providerserver.Serve(ctx, jtafProvider.NewProvider, opts)
	if err != nil {
		log.Fatal(err)
	}
}
