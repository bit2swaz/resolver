package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bit2swaz/resolver/internal/cache"
	"github.com/bit2swaz/resolver/internal/models"
	"github.com/bit2swaz/resolver/internal/scheduler"
)

const (
	dataDirPath   = "data"
	buildFilePath = "data/build.json"
	cacheFilePath = "data/cache.json"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "init":
		err = runInit()
	case "build":
		err = runBuild()
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runInit() error {
	if err := os.MkdirAll(dataDirPath, 0o755); err != nil {
		return err
	}

	targets := sampleTargets()
	if err := writeTargets(buildFilePath, targets); err != nil {
		return err
	}

	if err := cache.SaveState(cacheFilePath, &models.CacheState{Artifacts: map[string]string{}}); err != nil {
		return err
	}

	fmt.Printf("initialized sample files at %s and %s\n", buildFilePath, cacheFilePath)
	return nil

}

func runBuild() error {
	targets, err := loadTargets(buildFilePath)
	if err != nil {
		return err
	}

	if err := scheduler.Run(targets, cacheFilePath); err != nil {
		return err
	}

	for _, target := range targets {
		status := "executed"
		if target.IsCached {
			status = "cached"
		}

		fmt.Printf("%s: %s\n", target.ID, status)
	}

	return nil
}

func loadTargets(path string) ([]*models.Target, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var targets []*models.Target
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, err
	}

	return targets, nil
}

func writeTargets(path string, targets []*models.Target) error {
	data, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func sampleTargets() []*models.Target {
	return []*models.Target{
		{ID: "app", Dependencies: []string{"lib", "util"}, FileHash: "hash-app"},
		{ID: "lib", Dependencies: []string{"core"}, FileHash: "hash-lib"},
		{ID: "util", Dependencies: []string{"core"}, FileHash: "hash-util"},
		{ID: "core", Dependencies: nil, FileHash: "hash-core"},
	}
}

func printUsage() {
	program := filepath.Base(os.Args[0])
	fmt.Printf("usage: %s <init|build>\n", program)
}
