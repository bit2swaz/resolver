# Roadmap

This roadmap is now complete for the planned core implementation. The sections below preserve the phase structure while recording the concrete deliverables that exist in the repository today.

## Phase 1: Foundation and models

### Mini phase 1.1: module initialization and core schemas

- Status: completed
- Deliverables:
  - `go.mod` initialized with module path `github.com/bit2swaz/resolver`
  - `internal/models/models.go` defines `Target`, `Graph`, `CacheState`, and `AVLNode`
  - `internal/models/models_test.go` validates the shared model shapes
- Verification: `go test ./internal/models/...`

## Phase 2: The cache registry

### Mini phase 2.1: hash table implementation

- Status: completed
- Deliverables:
  - `internal/cache/hashtable.go` implements `HashTable`, `Set`, and `Get`
  - `internal/cache/hashtable_test.go` verifies insert and lookup behavior
- Verification: `go test ./internal/cache/...`

### Mini phase 2.2: AVL tree implementation

- Status: completed
- Deliverables:
  - `internal/cache/avltree.go` implements insertion, search, and balancing rotations
  - `internal/cache/avltree_test.go` verifies left rotation, right rotation, and search behavior
- Verification: `go test ./internal/cache/...`

## Phase 3: The dependency graph

### Mini phase 3.1: graph builder and cycle detection (DFS)

- Status: completed
- Deliverables:
  - `internal/graph/graph.go` builds adjacency lists from target definitions
  - DFS cycle detection reports circular dependencies as errors
  - `internal/graph/graph_test.go` covers both acyclic and cyclic graphs
- Verification: `go test ./internal/graph/...`

### Mini phase 3.2: topological sort scheduling

- Status: completed
- Deliverables:
  - `internal/graph/graph.go` exports `TopologicalSort`
  - DFS appends each node after its dependencies, producing a dependency-safe order
  - `internal/graph/graph_test.go` asserts dependencies appear before dependents
- Verification: `go test ./internal/graph/...`

## Phase 4: Persistence and execution

### Mini phase 4.1: JSON persistence manager

- Status: completed
- Deliverables:
  - `internal/cache/persistence.go` implements `SaveState` and `LoadState`
  - `internal/cache/persistence_test.go` verifies a JSON roundtrip using a temporary file
- Verification: `go test ./internal/cache/...`

### Mini phase 4.2: the build scheduler loop

- Status: completed
- Deliverables:
  - `internal/scheduler/scheduler.go` orchestrates cache loading, graph validation, sorting, execution decisions, and cache saving
  - The scheduler rebuilds in-memory hash table and AVL tree structures from persisted state
  - `internal/scheduler/scheduler_test.go` covers clean and fully cached runs
- Verification: `go test ./internal/scheduler/...`

## Phase 5: The command line interface

### Mini phase 5.1: main application wiring

- Status: completed
- Deliverables:
  - `cmd/resolver/main.go` implements `init` and `build` commands
  - `init` writes sample `data/build.json` and empty `data/cache.json`
  - `build` runs the scheduler and prints `executed` or `cached` per target
- Verification:
  - `go run ./cmd/resolver init`
  - `go run ./cmd/resolver build`

## Notes

- The current CLI exposes `init` and `build`. There is no separate `check` command in the implementation.
- The core roadmap is complete; future work can build on the existing scheduler, graph, and persistence layers.