// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/v66/github"
)

func ValidateRepositoryContent(content github.RepositoryContent) error {
	if content.DownloadURL == nil && content.Type != nil && *content.Type != "dir" {
		return fmt.Errorf("requested content has nil DownloadURL element")
	}

	if content.HTMLURL == nil {
		return fmt.Errorf("requested content has nil HTMLURL element")
	}

	if content.Path == nil {
		return fmt.Errorf("requested content has nil Path element")
	}

	if content.SHA == nil {
		return fmt.Errorf("requested content has nil SHA element")
	}

	_, err := url.Parse(*content.HTMLURL)
	if err != nil {
		return fmt.Errorf("while validating content HTMLURL element - %w", err)
	}

	return nil
}
