package graph

import (
	"testing"

	"github.com/bit2swaz/resolver/internal/models"
)

func TestBuildGraphAndHasCyclesWithAcyclicTargets(t *testing.T) {
	targets := []*models.Target{
		{ID: "app", Dependencies: []string{"lib", "util"}},
		{ID: "lib", Dependencies: []string{"core"}},
		{ID: "util", Dependencies: nil},
		{ID: "core", Dependencies: nil},
	}

	g := BuildGraph(targets)
	if g == nil {
		t.Fatal("expected graph to be built")
	}

	if len(g.Vertices) != len(targets) {
		t.Fatalf("expected %d vertices, got %d", len(targets), len(g.Vertices))
	}

	if len(g.Edges["app"]) != 2 {
		t.Fatalf("expected app to have 2 dependencies, got %d", len(g.Edges["app"]))
	}

	if err := HasCycles(g); err != nil {
		t.Fatalf("expected no cycle error, got %v", err)
	}
}

func TestHasCyclesReturnsErrorForCircularDependency(t *testing.T) {
	targets := []*models.Target{
		{ID: "A", Dependencies: []string{"B"}},
		{ID: "B", Dependencies: []string{"A"}},
	}

	g := BuildGraph(targets)
	err := HasCycles(g)
	if err == nil {
		t.Fatal("expected cycle detection error, got nil")
	}
}
