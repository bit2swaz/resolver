package models

// Target describes a build target and its dependency metadata.
type Target struct {
	ID           string
	Dependencies []string
	FileHash     string
	IsCached     bool
}

// Graph stores build targets and their adjacency list.
type Graph struct {
	Vertices []*Target
	Edges    map[string][]string
}

// CacheState stores cached artifact hashes for persistence.
type CacheState struct {
	Artifacts map[string]string
}

// AVLNode represents a single node in an AVL tree.
type AVLNode struct {
	Key    string
	Hash   string
	Height int
	Left   *AVLNode
	Right  *AVLNode
}
