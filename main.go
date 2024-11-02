package main

import (
	"context"
	"encoding/json"
	"fmt"
	ouryang "github.com/chrismarget-j/jtaf/yang"
	"log"
	"net/http"
	"strings"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
	"github.com/chrismarget-j/jtaf/junos"
	yangcache "github.com/chrismarget-j/jtaf/yang_cache"
	"github.com/openconfig/goyang/pkg/yang"
)

const configRoot = "junos-conf-root"

func main() {
	ctx := context.Background()

	cfg, err := jtafCfg.Get()
	if err != nil {
		log.Fatal(fmt.Errorf("while getting configuration - %w", err))
	}

	httpClient := http.DefaultClient

	yangDirs, err := yangcache.Populate(ctx, cfg, httpClient)
	if err != nil {
		log.Fatal(fmt.Errorf("while populating yang cache - %w", err))
	}

	breadcrumbTrails, err := junos.GetConfigBreadcrumbTrails(cfg.JunosConfigFile)
	if err != nil {
		log.Fatal(fmt.Errorf("while parsing device config file %q - %w", cfg.JunosConfigFile, err))
	}

	cfgRoot, err := ouryang.GetYangEntryByName(configRoot, yangDirs)
	if err != nil {
		log.Fatal(fmt.Errorf("while getting %s from %q - %w", configRoot, cfg.YangCacheDir(), err))
	}

	ouryang.TrimToConfig(cfgRoot, breadcrumbTrails)

	bytes, err := json.MarshalIndent(cfgRoot, "", "")
	if err != nil {
		log.Fatal(fmt.Errorf("while marshaling trimmed config root - %w", err))
	}

	fmt.Println(string(bytes))

	err = exploreYangBreadcrumbs(cfgRoot, breadcrumbTrails)
	if err != nil {
		log.Fatal(fmt.Errorf("while checking config breadcrumbs against yang data - %w", err))
	}

	getTypeKinds(cfgRoot)
	fmt.Println("kinds in use:\n", typeKinds)

	fmt.Println(cfgRoot.Name)
	for _, bct := range breadcrumbTrails {
		fmt.Println(strings.Join(bct, junos.PathSep))
	}
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

//var unions = make(map[string][]string)
//
//func getUnions(path string, ye *yang.Entry) {
//	if ye.Type != nil && ye.Type.Kind == yang.Yunion {
//		s := unions[path]
//	}
//}

func exploreYangBreadcrumbs(entry *yang.Entry, paths [][]string) error {
	for p, path := range paths {
		var ne *yang.Entry
		for s, step := range path {
			if ne == nil {
				ne = entry
			}

			var ok bool
			if ne, ok = ne.Dir[step]; !ok {
				return fmt.Errorf("failed traversing yang path %q at step %d (%s)", strings.Join(path, junos.PathSep), s, step)
			}
		}

		fmt.Printf("path %3d (%s) okay\n", p, strings.Join(path, junos.PathSep))
	}

	return nil
}
