package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
)

const xpathFile = "xpath_inputs.xml"

func main() {
	ctx := context.Background()

	cfg, err := getConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("while getting configuration - %w", err))
	}

	httpClient := http.DefaultClient

	err = populateYangCache(ctx, cfg, httpClient)
	if err != nil {
		log.Fatal(fmt.Errorf("while populating yang cache - %w", err))
	}

	xPaths, err := deviceConfigXpaths(cfg.DeviceConfigFile)
	if err != nil {
		log.Fatal(fmt.Errorf("while parsing user config Xpaths - %w", err))
	}

	b, err := xml.MarshalIndent(xPaths, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("while marshaling xpaths from device configuration - %w", err))
	}

	err = os.WriteFile(xpathFile, b, 0o644)
	if err != nil {
		log.Fatal(fmt.Errorf("while writing xpath data to %q - %w", xpathFile, err))
	}
}
