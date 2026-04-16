package graph

import (
	"strings"
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

func TestHasCyclesReturnsErrorForMissingDependency(t *testing.T) {
	targets := []*models.Target{
		{ID: "app", Dependencies: []string{"missing"}},
	}

	err := HasCycles(BuildGraph(targets))
	if err == nil {
		t.Fatal("expected missing dependency error, got nil")
	}

	if !strings.Contains(err.Error(), "missing dependency") {
		t.Fatalf("expected missing dependency error, got %v", err)
	}
}

func TestTopologicalSortPlacesDependenciesBeforeDependents(t *testing.T) {
	targets := []*models.Target{
		{ID: "app", Dependencies: []string{"lib", "util"}},
		{ID: "util", Dependencies: []string{"core"}},
		{ID: "lib", Dependencies: []string{"core"}},
		{ID: "core", Dependencies: nil},
	}

	g := BuildGraph(targets)
	order, err := TopologicalSort(g)
	if err != nil {
		t.Fatalf("expected topological sort to succeed, got %v", err)
	}

	if len(order) != len(targets) {
		t.Fatalf("expected %d items in order, got %d", len(targets), len(order))
	}

	positions := make(map[string]int, len(order))
	for index, id := range order {
		positions[id] = index
	}

	assertBefore(t, positions, "core", "lib")
	assertBefore(t, positions, "core", "util")
	assertBefore(t, positions, "lib", "app")
	assertBefore(t, positions, "util", "app")
}

func TestTopologicalSortHandlesDisconnectedGraphs(t *testing.T) {
	targets := []*models.Target{
		{ID: "frontend", Dependencies: []string{"ui-core"}},
		{ID: "ui-core", Dependencies: nil},
		{ID: "worker", Dependencies: []string{"queue"}},
		{ID: "queue", Dependencies: nil},
	}

	order, err := TopologicalSort(BuildGraph(targets))
	if err != nil {
		t.Fatalf("expected disconnected graph to sort, got %v", err)
	}

	if len(order) != len(targets) {
		t.Fatalf("expected %d items in order, got %d", len(targets), len(order))
	}

	positions := make(map[string]int, len(order))
	for index, id := range order {
		positions[id] = index
	}

	assertBefore(t, positions, "ui-core", "frontend")
	assertBefore(t, positions, "queue", "worker")
}

func assertBefore(t *testing.T, positions map[string]int, first, second string) {
	t.Helper()

	if positions[first] >= positions[second] {
		t.Fatalf("expected %s before %s, got positions %d and %d", first, second, positions[first], positions[second])
	}
}
