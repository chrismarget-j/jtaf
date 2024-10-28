package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

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

	//fmt.Print(strings.Join(xPaths, "\n") + "\n")
	b, err := xml.MarshalIndent(xPaths, "", "  ")
	fmt.Print(string(b) + "\n")
}
