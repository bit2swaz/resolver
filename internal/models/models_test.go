package models

import "testing"

func TestCoreStructsCanBeInstantiated(t *testing.T) {
	target := &Target{
		ID:           "app",
		Dependencies: []string{"lib"},
		FileHash:     "hash-app",
		IsCached:     false,
	}

	graph := &Graph{
		Vertices: []*Target{target},
		Edges: map[string][]string{
			"app": {"lib"},
		},
	}

	cacheState := &CacheState{
		Artifacts: map[string]string{
			"app": "hash-app",
		},
	}

	avlNode := &AVLNode{
		Key:    "app",
		Hash:   "hash-app",
		Height: 1,
	}

	if target.ID != "app" {
		t.Fatalf("expected target ID app, got %q", target.ID)
	}

	if len(graph.Vertices) != 1 {
		t.Fatalf("expected 1 graph vertex, got %d", len(graph.Vertices))
	}

	if cacheState.Artifacts["app"] != "hash-app" {
		t.Fatalf("expected cache artifact hash-app, got %q", cacheState.Artifacts["app"])
	}

	if avlNode.Key != "app" {
		t.Fatalf("expected avl node key app, got %q", avlNode.Key)
	}
}
