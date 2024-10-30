package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/openconfig/goyang/pkg/yang"
)

const xpathFile = "xpath_inputs.xml"

func main() {
	ctx := context.Background()

	cfg, err := getConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("while getting configuration - %w", err))
	}

	httpClient := http.DefaultClient

	yangDirs, err := populateYangCache(ctx, cfg, httpClient)
	if err != nil {
		log.Fatal(fmt.Errorf("while populating yang cache - %w", err))
	}

	//err = deviceConfigXpathsToFile(cfg.DeviceConfigFile, xpathFile)
	//if err != nil {
	//	log.Fatal(fmt.Errorf("while parsing device configs and writing to xpath file - %w", err))
	//}

	// deviceConfigXpaths()

	breadcrumbTrails, err := getConfigBreadcrumbTrails(cfg.DeviceConfigFile)
	if err != nil {
		log.Fatal(fmt.Errorf("while parsing device config file %q - %w", cfg.DeviceConfigFile, err))
	}

	cfgRoot, err := getYangEntryConfigRoot(yangDirs)
	if err != nil {
		log.Fatal(fmt.Errorf("while getting %s from %q - %w", yangConfigRoot, cfg.yangCacheDir(), err))
	}

	err = yangWalk(cfgRoot, breadcrumbTrails)
	if err != nil {
		log.Fatal(fmt.Errorf("while checking config breadcrumbs against yang data - %w", err))
	}

	getTypeKinds(cfgRoot)
	fmt.Println("kinds in use:\n", typeKinds)

	//fmt.Println(cfgRoot.Name)
	//for _, bct := range breadcrumbTrails {
	//	fmt.Println(strings.Join(bct, pathSep))
	//}
}

var typeKinds = make(map[yang.TypeKind]struct{})

func getTypeKinds(ye *yang.Entry) {
	if ye.Type != nil {
		typeKinds[ye.Type.Kind] = struct{}{}
	}

	for _, ye := range ye.Dir {
		getTypeKinds(ye)
	}
}
