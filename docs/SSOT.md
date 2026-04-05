# ssot: distributed build cache dependency resolver

## 1. project overview
* **objective:** build a dependency resolver and task scheduler that models software build targets as a directed acyclic graph (DAG), resolves the correct execution order, and uses caching to skip redundant tasks.
* **language:** Go
* **primary comparisons:** cache lookups using a Hash Table vs an AVL Tree, and dependency traversal using Depth First Search vs Breadth First Search.

## 2. core data structures and syllabus mapping
* **dependency graph (graphs):** represents the build targets as an adjacency list. vertices are stored in an array, and edges are stored as linked lists attached to each vertex.
* **cache registry (trees and hash tables):** stores metadata of compiled artifacts. the primary implementation is a Hash Table for O(1) lookups. the secondary implementation is an AVL Tree for O(log n) lookups to satisfy the comparison requirement.
* **execution queues (queues and stacks):** a FIFO queue is used to hold tasks that are ready to build. a LIFO stack is used implicitly or explicitly during the Depth First Search cycle detection phase.

## 3. algorithmic workflow
* **initialization:** load the previous cache state and dependency definitions from a local JSON file.
* **validation:** run Depth First Search on the graph to detect cycles. if a back-edge is found, throw an explicit error and halt execution.
* **topological sort:** use Kahn's algorithm with the FIFO queue to flatten the DAG into a linear, executable build order.
* **cache check:** query the cache registry for each target in the sorted list. if the target is unchanged and exists in the cache, mark it as cached.
* **execution:** process uncached targets, simulate a build step, and insert their new hashes into both the Hash Table and the AVL Tree.
* **teardown:** serialize and save the updated cache state back to the JSON file.

## 4. validation strategy (the 6 test cases)
* **clean build:** an empty cache where all targets must be executed and stored.
* **fully cached build:** no source files have changed, meaning the system resolves the graph and bypasses execution by hitting the cache.
* **partial rebuild:** one leaf node changes, triggering a rebuild of only that node and its direct dependents.
* **disconnected graph:** a configuration with two entirely separate build trees that schedule and build independently.
* **missing file error:** a target references a dependency that does not exist, triggering a system error.
* **circular dependency:** an intentional cycle is introduced to ensure the validation catches it and aborts.

## 5. report comparison points
* **time complexity:** mathematical proofs and benchmarks showing lookup speed differences between the Hash Table and the AVL Tree.
* **space complexity:** memory overhead comparison between the adjacency list and the cache registry structures.
* **execution fairness:** explanation of how a FIFO queue appropriately handles topological sorting compared to a LIFO stack.

## 6. directory structure
* **cmd/resolver/main.go:** entry point for the command line interface.
* **internal/graph/:** contains the adjacency list, topological sort, and cycle detection logic.
* **internal/cache/:** contains the Hash Table and AVL Tree implementations, plus the JSON persistence logic.
* **internal/scheduler/:** contains the FIFO queue and the main build execution loop.
* **data/:** directory to store the config and state JSON files.

## 7. core data schemas (Go structs)
* **Target:** struct containing ID (string), Dependencies (slice of strings), FileHash (string), and IsCached (boolean).
* **Graph:** struct containing Vertices (slice of Target pointers) and Edges (map of string to slice of strings, representing the adjacency list).
* **CacheState:** struct containing Artifacts (map of string to string) for JSON serialization.
* **AVLNode:** struct containing Key (string), Hash (string), Height (int), Left (pointer to AVLNode), and Right (pointer to AVLNode).

## 8. command line interface flow
* **go run cmd/resolver/main.go init:** creates the initial empty state file and a sample build configuration file.
* **go run cmd/resolver/main.go check:** parses the graph and runs the cycle detection without building anything.
* **go run cmd/resolver/main.go build:** runs the full pipeline for validation, sorting, cache checking, and execution.