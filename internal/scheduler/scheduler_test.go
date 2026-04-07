package scheduler

import (
	"path/filepath"
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

func sampleTargets() []*models.Target {
	return []*models.Target{
		{ID: "app", Dependencies: []string{"lib", "util"}, FileHash: "hash-app"},
		{ID: "lib", Dependencies: []string{"core"}, FileHash: "hash-lib"},
		{ID: "util", Dependencies: []string{"core"}, FileHash: "hash-util"},
		{ID: "core", Dependencies: nil, FileHash: "hash-core"},
	}
}
