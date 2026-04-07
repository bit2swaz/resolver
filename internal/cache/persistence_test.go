package cache

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bit2swaz/resolver/internal/models"
)

func TestSaveStateAndLoadStateRoundTrip(t *testing.T) {
	want := &models.CacheState{
		Artifacts: map[string]string{
			"app":  "hash-app",
			"lib":  "hash-lib",
			"util": "hash-util",
		},
	}

	statePath := filepath.Join(t.TempDir(), "cache.json")

	if err := SaveState(statePath, want); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	got, err := LoadState(statePath)
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected loaded state %#v, got %#v", want, got)
	}
}
