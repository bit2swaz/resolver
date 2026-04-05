package cache

import "testing"

func TestHashTableStoresAndRetrievesHashes(t *testing.T) {
	table := NewHashTable()
	table.Set("app", "hash-app")

	hash, ok := table.Get("app")
	if !ok {
		t.Fatal("expected hash for app to exist")
	}

	if hash != "hash-app" {
		t.Fatalf("expected hash-app, got %q", hash)
	}
}

func TestHashTableReturnsFalseForMissingKeys(t *testing.T) {
	table := NewHashTable()

	_, ok := table.Get("missing")
	if ok {
		t.Fatal("expected missing key lookup to return false")
	}
}
