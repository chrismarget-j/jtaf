package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	err = deviceConfigXpathsToFile(cfg.DeviceConfigFile, xpathFile)
	if err != nil {
		log.Fatal(fmt.Errorf("while parsing device configs and writing to xpath file - %w", err))
	}
}
