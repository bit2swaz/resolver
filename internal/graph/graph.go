package graph

import (
	"fmt"

	"github.com/bit2swaz/resolver/internal/models"
)

// BuildGraph converts targets into the shared graph model with adjacency lists.
func BuildGraph(targets []*models.Target) *models.Graph {
	edges := make(map[string][]string, len(targets))

	for _, target := range targets {
		dependencies := append([]string(nil), target.Dependencies...)
		edges[target.ID] = dependencies
	}

	return &models.Graph{
		Vertices: targets,
		Edges:    edges,
	}
}

// HasCycles runs DFS and returns an error when it finds a back-edge.
func HasCycles(g *models.Graph) error {
	if err := validateGraph(g); err != nil {
		return err
	}

	visited := make(map[string]bool, len(g.Vertices))
	onPath := make(map[string]bool, len(g.Vertices))
	pathStack := make([]string, 0, len(g.Vertices))

	var dfs func(string) error
	dfs = func(node string) error {
		visited[node] = true
		onPath[node] = true
		pathStack = append(pathStack, node)

		for _, dependency := range g.Edges[node] {
			if onPath[dependency] {
				return fmt.Errorf("cycle detected involving %s", dependency)
			}

			if visited[dependency] {
				continue
			}

			if err := dfs(dependency); err != nil {
				return err
			}
		}

		onPath[node] = false
		pathStack = pathStack[:len(pathStack)-1]
		return nil
	}

	for _, vertex := range g.Vertices {
		if visited[vertex.ID] {
			continue
		}

		if err := dfs(vertex.ID); err != nil {
			return err
		}
	}

	return nil
}

// TopologicalSort returns a dependency-safe execution order using DFS.
func TopologicalSort(g *models.Graph) ([]string, error) {
	if err := validateGraph(g); err != nil {
		return nil, err
	}

	visited := make(map[string]bool, len(g.Vertices))
	onPath := make(map[string]bool, len(g.Vertices))
	stack := make([]string, 0, len(g.Vertices))

	var dfs func(string) error
	dfs = func(node string) error {
		visited[node] = true
		onPath[node] = true

		for _, dependency := range g.Edges[node] {
			if onPath[dependency] {
				return fmt.Errorf("cycle detected involving %s", dependency)
			}

			if visited[dependency] {
				continue
			}

			if err := dfs(dependency); err != nil {
				return err
			}
		}

		onPath[node] = false
		stack = append(stack, node)
		return nil
	}

	for _, vertex := range g.Vertices {
		if visited[vertex.ID] {
			continue
		}

		if err := dfs(vertex.ID); err != nil {
			return nil, err
		}
	}

	order := make([]string, len(stack))
	copy(order, stack)
	return order, nil
}

func validateGraph(g *models.Graph) error {
	if g == nil {
		return fmt.Errorf("graph is nil")
	}

	knownTargets := make(map[string]struct{}, len(g.Vertices))
	for _, vertex := range g.Vertices {
		if vertex == nil {
			return fmt.Errorf("graph contains nil target")
		}

		if _, exists := knownTargets[vertex.ID]; exists {
			return fmt.Errorf("duplicate target definition for %s", vertex.ID)
		}

		knownTargets[vertex.ID] = struct{}{}
	}

	for _, vertex := range g.Vertices {
		for _, dependency := range g.Edges[vertex.ID] {
			if _, exists := knownTargets[dependency]; !exists {
				return fmt.Errorf("missing dependency %s referenced by %s", dependency, vertex.ID)
			}
		}
	}

	return nil
}
