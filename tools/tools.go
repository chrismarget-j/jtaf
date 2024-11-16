// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build tools

package tools

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"

	// license compliance
	_ "github.com/google/go-licenses"

	// staticcheck
	_ "honnef.co/go/tools/cmd/staticcheck"

	// release
	_ "github.com/goreleaser/goreleaser"

	// gofumpt does strict formatting
	_ "mvdan.cc/gofumpt"
)
