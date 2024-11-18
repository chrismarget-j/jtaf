package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
	yangcache "github.com/chrismarget-j/jtaf/yang_cache"
)

func main() {
	ctx := context.Background()

	log.Printf("Reading configuration...")
	cfg, err := jtafCfg.Get()
	if err != nil {
		log.Fatal(fmt.Errorf("while getting configuration - %w", err))
	}

	httpClient := http.DefaultClient

	yangDirs, err := yangcache.Populate(ctx, cfg, httpClient) // 1.9s (validate) / 4-13s (from github)
	if err != nil {
		log.Fatal(fmt.Errorf("while populating yang cache - %w", err))
	}

	log.Println("yang files are in:")
	for _, d := range yangDirs {
		log.Println(" - " + d)
	}

}
