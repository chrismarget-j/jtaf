// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package jtafCfg

import "os"

const envPrefix = "JTAF_"

func fromEnv(name string, target *string) {
	if *target == "" {
		*target = os.Getenv(envPrefix + name)
	}
}
