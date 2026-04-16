package scheduler

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/bit2swaz/resolver/internal/cache"
	"github.com/bit2swaz/resolver/internal/models"
)

func TestRunPerformsCleanBuildAndPersistsCache(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := cache.SaveState(cachePath, &models.CacheState{Artifacts: map[string]string{}}); err != nil {
		t.Fatalf("expected initial cache save to succeed, got %v", err)
	}

	targets := sampleTargets()

	if err := Run(targets, cachePath); err != nil {
		t.Fatalf("expected clean build to succeed, got %v", err)
	}

	for _, target := range targets {
		if target.IsCached {
			t.Fatalf("expected %s to execute on clean build", target.ID)
		}
	}

	state, err := cache.LoadState(cachePath)
	if err != nil {
		t.Fatalf("expected saved cache to load, got %v", err)
	}

	if len(state.Artifacts) != len(targets) {
		t.Fatalf("expected %d cached artifacts, got %d", len(targets), len(state.Artifacts))
	}

	for _, target := range targets {
		if state.Artifacts[target.ID] != target.FileHash {
			t.Fatalf("expected cached hash %q for %s, got %q", target.FileHash, target.ID, state.Artifacts[target.ID])
		}
	}
}

func TestRunMarksTargetsCachedWhenHashesAlreadyExist(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	firstRunTargets := sampleTargets()

	if err := cache.SaveState(cachePath, &models.CacheState{Artifacts: map[string]string{}}); err != nil {
		t.Fatalf("expected initial cache save to succeed, got %v", err)
	}

	if err := Run(firstRunTargets, cachePath); err != nil {
		t.Fatalf("expected first run to succeed, got %v", err)
	}

	secondRunTargets := sampleTargets()
	if err := Run(secondRunTargets, cachePath); err != nil {
		t.Fatalf("expected second run to succeed, got %v", err)
	}

	for _, target := range secondRunTargets {
		if !target.IsCached {
			t.Fatalf("expected %s to be treated as cached on second run", target.ID)
		}
	}
}

func TestRunRebuildsChangedTargetAndItsDependents(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := cache.SaveState(cachePath, &models.CacheState{Artifacts: map[string]string{}}); err != nil {
		t.Fatalf("expected initial cache save to succeed, got %v", err)
	}

	if err := Run(sampleTargets(), cachePath); err != nil {
		t.Fatalf("expected first run to succeed, got %v", err)
	}

	changedTargets := sampleTargets()
	for _, target := range changedTargets {
		if target.ID == "util" {
			target.FileHash = "hash-util-v2"
		}
	}

	if err := Run(changedTargets, cachePath); err != nil {
		t.Fatalf("expected partial rebuild run to succeed, got %v", err)
	}

	statuses := make(map[string]bool, len(changedTargets))
	for _, target := range changedTargets {
		statuses[target.ID] = target.IsCached
	}

	if !statuses["core"] {
		t.Fatal("expected core to remain cached when util changes")
	}

	if !statuses["lib"] {
		t.Fatal("expected lib to remain cached when util changes")
	}

	if statuses["util"] {
		t.Fatal("expected util to rebuild after its hash changes")
	}

	if statuses["app"] {
		t.Fatal("expected app to rebuild when a dependency rebuilds")
	}

	state, err := cache.LoadState(cachePath)
	if err != nil {
		t.Fatalf("expected updated cache state to load, got %v", err)
	}

	if state.Artifacts["util"] != "hash-util-v2" {
		t.Fatalf("expected updated util hash to persist, got %q", state.Artifacts["util"])
	}
}

func TestRunReturnsErrorForMissingDependency(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := cache.SaveState(cachePath, &models.CacheState{Artifacts: map[string]string{}}); err != nil {
		t.Fatalf("expected initial cache save to succeed, got %v", err)
	}

	targets := []*models.Target{
		{ID: "app", Dependencies: []string{"missing"}, FileHash: "hash-app"},
	}

	err := Run(targets, cachePath)
	if err == nil {
		t.Fatal("expected missing dependency error, got nil")
	}

	if !strings.Contains(err.Error(), "missing dependency") {
		t.Fatalf("expected missing dependency error, got %v", err)
	}
}

func sampleTargets() []*models.Target {
	return []*models.Target{
		{ID: "app", Dependencies: []string{"lib", "util"}, FileHash: "hash-app"},
		{ID: "lib", Dependencies: []string{"core"}, FileHash: "hash-lib"},
		{ID: "util", Dependencies: []string{"core"}, FileHash: "hash-util"},
		{ID: "core", Dependencies: nil, FileHash: "hash-core"},
	}
}
