# Resolver

Resolver is a Go command-line project that models build targets as a dependency graph, computes a dependency-safe execution order with depth-first search, and persists a simple build cache between runs.

## What it does

- Loads build targets from `data/build.json`
- Loads cached artifact hashes from `data/cache.json`
- Builds an adjacency-list graph from the targets
- Detects circular dependencies with DFS
- Rejects missing dependency references before scheduling begins
- Produces a dependency-safe execution order with DFS-based topological sorting
- Checks a hash table cache before simulating execution
- Rebuilds a changed target and its downstream dependents in the same run
- Updates both the persisted cache state and an AVL tree used for the project comparison requirement

## Project layout

- `cmd/resolver/main.go`: CLI entrypoint with `init` and `build` commands
- `internal/models`: shared data models
- `internal/cache`: hash table, AVL tree, and JSON persistence helpers
- `internal/graph`: graph construction, cycle detection, and topological sort
- `internal/scheduler`: orchestration of graph validation, ordering, cache checks, and cache updates

## Commands

### Initialize sample data

```bash
go run ./cmd/resolver init
```

Creates the `data/` directory along with:

- `data/build.json`: sample targets (`app`, `lib`, `util`, `core`)
- `data/cache.json`: an empty cache state

### Run the build pipeline

```bash
go run ./cmd/resolver build
```

This command reads the JSON files, runs the scheduler, and prints each target as either `executed` or `cached`.

Run `init` first so `data/cache.json` exists before the scheduler starts.

## Data model

Each build target uses the shared `Target` schema:

```json
{
	"ID": "app",
	"Dependencies": ["lib", "util"],
	"FileHash": "hash-app",
	"IsCached": false
}
```

The persisted cache file stores artifact hashes by target ID:

```json
{
	"Artifacts": {
		"core": "hash-core"
	}
}
```

## Testing

The implementation is covered by package-level tests for:

- core models
- hash table cache operations
- AVL tree balancing and lookup
- JSON persistence
- graph construction, disconnected graph ordering, and cycle detection
- missing dependency validation
- DFS-based topological ordering
- scheduler behavior for clean, fully cached, and partial rebuild runs

Benchmark-based comparison data for the hash table and AVL tree implementations is available in `docs/PROJECT_REPORT.md`.

Run the test suite with:

```bash
go test ./...
```

## Current status

The core roadmap phases are implemented:

- foundation and shared models
- hash table and AVL tree cache structures
- dependency graph construction and cycle detection
- topological sorting
- JSON cache persistence
- scheduler loop
- CLI wiring

See `docs/SSOT.md` for the current design summary, `docs/ROADMAP.md` for implementation status, and `docs/PROJECT_REPORT.md` for the full mini-project report.