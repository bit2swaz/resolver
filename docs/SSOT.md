# SSOT: Resolver design summary

## 1. Project overview

- Objective: model build targets as a directed acyclic dependency graph, compute a safe execution order, and skip redundant work through persisted caching.
- Language: Go
- Implemented comparison structures: hash table for primary cache lookup and AVL tree for the secondary academic comparison requirement.

## 2. Core data structures

### Target

- Fields:
  - `ID string`
  - `Dependencies []string`
  - `FileHash string`
  - `IsCached bool`
- Purpose: describes one build target and its dependencies.

### Graph

- Fields:
  - `Vertices []*Target`
  - `Edges map[string][]string`
- Purpose: stores the target list and adjacency lists used by graph traversal.

### CacheState

- Fields:
  - `Artifacts map[string]string`
- Purpose: persisted JSON cache keyed by target ID.

### AVLNode

- Fields:
  - `Key string`
  - `Hash string`
  - `Height int`
  - `Left *AVLNode`
  - `Right *AVLNode`
- Purpose: in-memory AVL tree node used for the comparison implementation.

## 3. Package layout

- `cmd/resolver/main.go`: application entrypoint
- `internal/models`: shared structs
- `internal/cache`: hash table, AVL tree, and persistence helpers
- `internal/graph`: graph construction, cycle detection, and topological sort
- `internal/scheduler`: end-to-end scheduler pipeline

## 4. Runtime workflow

1. Load targets from `data/build.json`.
2. Load the existing cache state from `data/cache.json`.
3. Rebuild in-memory hash table and AVL tree structures from persisted cache data.
4. Build a graph from the configured targets.
5. Run DFS cycle detection and stop if a circular dependency is found.
6. Run DFS-based topological sorting to produce a dependency-safe order.
7. Iterate over the ordered targets, checking the hash table cache first.
8. Mark matching targets as cached; otherwise simulate execution and update the cache state.
9. Save the updated cache state back to `data/cache.json`.

## 5. Implemented algorithm details

### Cycle detection

- Uses depth-first search.
- Tracks visited nodes and the current recursion path.
- Reports an error when DFS encounters a dependency already on the current path.

### Topological sorting

- Uses depth-first search.
- Appends a node only after visiting its dependencies.
- Returns a dependency-safe order where dependencies appear before dependents.

### Cache lookup

- Primary lookup path: hash table for constant-time reads.
- Secondary structure: AVL tree updated alongside the hash table for comparison purposes.

## 6. CLI behavior

### `init`

- Creates the `data/` directory.
- Writes a sample `data/build.json` file.
- Writes an empty `data/cache.json` file.

### `build`

- Loads targets and cache state from disk.
- Runs the scheduler.
- Prints each target with either `executed` or `cached`.

## 7. Validation coverage

The current automated tests cover:

- model instantiation
- hash table insertion and lookup
- AVL rotations and search
- JSON persistence roundtrips
- acyclic and cyclic graph behavior
- dependency ordering in topological sort
- clean and fully cached scheduler runs

## 8. Current scope notes

- The implemented CLI includes `init` and `build` only.
- The scheduler expects the cache file to exist, so the intended first-run flow is `init` followed by `build`.
- This repository currently documents and ships the completed core implementation rather than future aspirational phases.