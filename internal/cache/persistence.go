package cache

import (
	"encoding/json"
	"os"

	"github.com/bit2swaz/resolver/internal/models"
)

// SaveState writes the cache state to disk as JSON.
func SaveState(path string, state *models.CacheState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// LoadState reads the cache state from disk and decodes it from JSON.
func LoadState(path string) (*models.CacheState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	state := &models.CacheState{}
	if err := json.Unmarshal(data, state); err != nil {
		return nil, err
	}

	if state.Artifacts == nil {
		state.Artifacts = make(map[string]string)
	}

	return state, nil
}
