// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
	"github.com/chrismarget-j/jtaf/junos"
	ouryang "github.com/chrismarget-j/jtaf/yang"
	yangcache "github.com/chrismarget-j/jtaf/yang_cache"
	"github.com/openconfig/goyang/pkg/yang"
)

const (
	junosByRefPattern = "<.*>|$.*"
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

	breadcrumbTrails, err := junos.GetConfigBreadcrumbTrails(cfg.JunosConfigFile)
	if err != nil {
		log.Fatal(fmt.Errorf("while parsing device config file %q - %w", cfg.JunosConfigFile, err))
	}

	cfgRoot, err := ouryang.GetYangEntryByName(cfg.ConfigRoot(), yangDirs) // 4.5s
	if err != nil {
		log.Fatal(fmt.Errorf("while getting %s from %q - %w", cfg.ConfigRoot(), yangDirs, err))
	}

	log.Printf("Trimming YANG entries down to the minimal set required for %d configuration paths...\n", len(breadcrumbTrails))
	ouryang.TrimToConfig(cfgRoot, breadcrumbTrails) // <1ms

	ouryang.TrimUnions(cfgRoot) // < 1ms
	return

	getUnions(cfgRoot, nil)
	ue := unionEntries

	us := make(map[string]string)
entryLoop:
	for k, v := range ue {
		if len(v.Type.Type) == 2 {
			for _, t := range v.Type.Type {
				if t.Kind == yang.Ystring && len(t.Pattern) == 1 && t.Pattern[0] == junosByRefPattern {
					continue entryLoop
				}
			}
		}
		s, err := ouryang.TypeToString(v.Type)
		if err != nil {
			log.Fatal(err)
		}
		us[k] = s
		//if len(v.Type) ==  2 [
		//
		//]
		fmt.Printf("%s: %s\n", k, s)
	}

	bytes, err := json.MarshalIndent(cfgRoot, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("while marshaling trimmed config root - %w", err))
	}

	_ = bytes
	// fmt.Println(string(bytes))

	err = exploreYangBreadcrumbs(cfgRoot, breadcrumbTrails)
	if err != nil {
		log.Fatal(fmt.Errorf("while checking config breadcrumbs against yang data - %w", err))
	}

	// getTypeKinds(cfgRoot)
	// fmt.Println("kinds in use:\n", typeKinds)

	//fmt.Println(cfgRoot.Name)
	//for _, bct := range breadcrumbTrails {
	//	fmt.Println(strings.Join(bct, junos.PathSep))
	//}
}

var unionEntries = make(map[string]*yang.Entry)

func getUnions(ye *yang.Entry, p []string) {
	for _, de := range ye.Dir {
		if de.Type != nil && de.Type.Kind == yang.Yunion {
			unionEntries[path.Join(append(p, de.Name)...)] = de
		}

		getUnions(de, append(p, de.Name))
	}
}

//var typeKinds = make(map[yang.TypeKind]struct{})
//
//func getTypeKinds(ye *yang.Entry) {
//	if ye.Type != nil {
//		typeKinds[ye.Type.Kind] = struct{}{}
//	}
//
//	for _, ye := range ye.Dir {
//		getTypeKinds(ye)
//	}
//}

//var unionEntries = make(map[string][]string)
//
//func getUnions(path string, ye *yang.Entry) {
//	if ye.Type != nil && ye.Type.Kind == yang.Yunion {
//		s := unionEntries[path]
//	}
//}

func exploreYangBreadcrumbs(entry *yang.Entry, paths [][]string) error {
	log.Println("Inspecting trimmed YANG entries...")
	for _, path := range paths {
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
	}

	return nil
}

func getLeafNames(ye *yang.Entry) map[string]struct{} {
	//var totalCount int                  // 6173
	//var alphaNumCount int               // 949
	//var alphaNumWithHyphenCount int     // 5210
	//var alphaNumWithUnderscoreCount int // 10
	//var remainingCount int
	//alphaNum := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")
	//alphaNumWithHyphen := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*$")
	//alphaNumWithUnderscore := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")
	//
	//leafNameMap := getLeafNames(cfgRoot)
	//totalCount = len(leafNameMap)
	//leafNameSlice := make([]string, 0, len(leafNameMap))
	//for k := range leafNameMap {
	//	if alphaNum.MatchString(k) {
	//		alphaNumCount++
	//		continue
	//	}
	//	if alphaNumWithHyphen.MatchString(k) {
	//		alphaNumWithHyphenCount++
	//		continue
	//	}
	//	if alphaNumWithUnderscore.MatchString(k) {
	//		alphaNumWithUnderscoreCount++
	//		continue
	//	}
	//	leafNameSlice = append(leafNameSlice, k)
	//}
	//sort.Strings(leafNameSlice)
	//remainingCount = len(leafNameSlice)
	//_ = totalCount
	//_ = remainingCount
	//
	//for _, ln := range leafNameSlice {
	//	fmt.Println(ln)
	//}

	result := make(map[string]struct{})
	if ye.Kind == yang.LeafEntry {
		result[ye.Name] = struct{}{}
	}

	for _, de := range ye.Dir {
		deResult := getLeafNames(de)
		for k := range deResult {
			result[k] = struct{}{}
		}
	}

	return result
}
