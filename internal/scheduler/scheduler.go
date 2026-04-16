package scheduler

import (
	"fmt"

	"github.com/bit2swaz/resolver/internal/cache"
	"github.com/bit2swaz/resolver/internal/graph"
	"github.com/bit2swaz/resolver/internal/models"
)

// Run executes the dependency resolution and cache update pipeline.
func Run(targets []*models.Target, cachePath string) error {
	state, err := cache.LoadState(cachePath)
	if err != nil {
		return err
	}

	hashTable := cache.NewHashTable()
	var avlRoot *models.AVLNode
	for id, hash := range state.Artifacts {
		hashTable.Set(id, hash)
		avlRoot = cache.Insert(avlRoot, id, hash)
	}

	g := graph.BuildGraph(targets)
	if err := graph.HasCycles(g); err != nil {
		return err
	}

	order, err := graph.TopologicalSort(g)
	if err != nil {
		return err
	}

	targetByID := make(map[string]*models.Target, len(targets))
	for _, target := range targets {
		targetByID[target.ID] = target
	}

	rebuiltTargets := make(map[string]bool, len(targets))
	queue := append([]string(nil), order...)
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]

		target, ok := targetByID[id]
		if !ok {
			return fmt.Errorf("missing target definition for %s", id)
		}

		dependencyRebuilt := false
		for _, dependency := range target.Dependencies {
			if rebuiltTargets[dependency] {
				dependencyRebuilt = true
				break
			}
		}

		cachedHash, ok := hashTable.Get(id)
		if !dependencyRebuilt && ok && cachedHash == target.FileHash {
			target.IsCached = true
			rebuiltTargets[id] = false
			continue
		}

		target.IsCached = false
		hashTable.Set(id, target.FileHash)
		avlRoot = cache.Insert(avlRoot, id, target.FileHash)
		state.Artifacts[id] = target.FileHash
		rebuiltTargets[id] = true
	}

	return cache.SaveState(cachePath, state)
}
