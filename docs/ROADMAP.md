### phase 1: foundation and models

**mini phase 1.1: module initialization and core schemas**

  * **purpose:** set up the go module and define the exact structs from the SSOT so all future packages have a single source of truth for data shapes.
  * **definition of done:** `go.mod` exists, and `internal/models/models.go` contains the definitions for Target, Graph, CacheState, and AVLNode. tests successfully instantiate these structs.
  * **verification:** `go test ./internal/models/...` passes.
  * **master prompt:**

> initialize a new go module named github.com/bit2swaz/resolver. use the red green refactor tdd method. first, write a failing test in `internal/models/models_test.go` that attempts to instantiate the following structs: Target (ID string, Dependencies []string, FileHash string, IsCached bool), Graph (Vertices []\*Target, Edges map[string][]string), CacheState (Artifacts map[string]string), and AVLNode (Key string, Hash string, Height int, Left \*AVLNode, Right \*AVLNode). then, create `internal/models/models.go` and implement these structs to make the test pass. keep the code clean and simple. add simple inline comments explaining what each struct does. run the tests. if they pass, commit and push the code with the message "feat: initialize go module and core data schemas". use context7

### phase 2: the cache registry

**mini phase 2.1: hash table implementation**

  * **purpose:** implement the primary $O(1)$ cache registry.
  * **definition of done:** `internal/cache/hashtable.go` exists with `Set` and `Get` methods.
  * **verification:** unit tests verify that items can be stored and retrieved accurately.
  * **master prompt:**

> following the red green refactor tdd method, create a new package `internal/cache`. first, write failing tests in `hashtable_test.go` for a HashTable struct that tests storing a file hash string by a target ID string, and retrieving it. then, implement the HashTable struct and its `Set(id, hash string)` and `Get(id string) (string, bool)` methods in `hashtable.go` using a standard go map to make the tests pass. keep the code clean and simple with simple inline comments. run the tests. if they pass, commit and push the code with the message "feat: implement hash table cache registry". use context7

**mini phase 2.2: AVL tree implementation**

  * **purpose:** implement the secondary $O(\log n)$ cache registry for your academic comparison requirement.
  * **definition of done:** `internal/cache/avltree.go` exists with standard insert, search, and self-balancing rotation methods.
  * **verification:** unit tests verify right and left rotations maintain tree balance, and search returns correct nodes.
  * **master prompt:**

> following the red green refactor tdd method, write failing tests in `internal/cache/avltree_test.go` that test inserting unordered target IDs and file hashes into an AVL tree, ensuring the tree remains balanced, and searching for an ID successfully retrieves the hash. use the AVLNode struct from the models package. then, implement `avltree.go` with `Insert(root *AVLNode, key, hash string) *AVLNode` and `Search(root *AVLNode, key string) *AVLNode` methods, including the necessary right and left rotation helper functions to maintain balance. make the tests pass. keep the code clean and simple with simple inline comments. run the tests. if they pass, commit and push the code with the message "feat: implement avl tree cache registry for academic comparison". use context7

### phase 3: the dependency graph

**mini phase 3.1: graph builder and cycle detection (DFS)**

  * **purpose:** parse targets into an adjacency list and implement DFS to catch circular dependencies (test case 6).
  * **definition of done:** `internal/graph/graph.go` can build the graph from a list of targets and accurately throw an error if a cycle exists.
  * **verification:** unit tests pass for a valid acyclic graph, and explicitly fail (returning an error) for a graph with a circular dependency.
  * **master prompt:**

> following the red green refactor tdd method, create `internal/graph/graph_test.go`. write two tests: one that builds a valid Graph from a list of Target structs and expects no cycles, and one that builds a Graph with an intentional circular dependency (target A depends on B, B depends on A) and expects a cycle detection error. then, implement `graph.go` with a `BuildGraph(targets []*Target) *Graph` function, and a `HasCycles(g *Graph) error` function that uses depth first search (DFS) and a LIFO stack approach to detect back-edges. make the tests pass. keep the code clean and simple with simple inline comments. run the tests. if they pass, commit and push the code with the message "feat: implement graph builder and dfs cycle detection". use context7

**mini phase 3.2: topological sort scheduling (Kahn's)**

  * **purpose:** flatten the valid DAG into an executable build order using a FIFO queue.
  * **definition of done:** the graph package exports a sort function that returns an ordered slice of target IDs.
  * **verification:** unit tests prove that parent dependencies are always sorted before the targets that rely on them.
  * **master prompt:**

> following the red green refactor tdd method, add a failing test in `internal/graph/graph_test.go` that provides a valid, unsorted Graph of targets and asserts that the resulting ordered array correctly places all dependencies before their dependents. then, implement a `TopologicalSort(g *Graph) ([]string, error)` function in `graph.go` using Kahn's algorithm. you must use a standard slice acting as a FIFO queue to process nodes with zero in-degree. make the test pass. keep the code clean and simple with simple inline comments. run the tests. if they pass, commit and push the code with the message "feat: implement topologial sorting using kahn's algorithm". use context7

### phase 4: persistence and execution

**mini phase 4.1: JSON persistence manager**

  * **purpose:** save and load the cache state and build configurations to disk.
  * **definition of done:** `internal/cache/persistence.go` can serialize CacheState to a file and read it back into memory.
  * **verification:** unit tests write dummy data to a temporary file, read it back, and assert the data matches exactly.
  * **master prompt:**

> following the red green refactor tdd method, create a test in `internal/cache/persistence_test.go` that creates a dummy CacheState struct, writes it to a temporary JSON file, reads it back into a new struct, and asserts they are equal. then, implement `persistence.go` with two functions: `SaveState(path string, state *models.CacheState) error` and `LoadState(path string) (*models.CacheState, error)`. make the tests pass. keep the code clean and simple with simple inline comments. run the tests. if they pass, commit and push the code with the message "feat: implement json state persistence". use context7

**mini phase 4.2: the build scheduler loop**

  * **purpose:** the central controller that orchestrates the graph validation, sorting, cache checking, and mock execution.
  * **definition of done:** `internal/scheduler/scheduler.go` runs the full pipeline and updates the cache.
  * **verification:** unit tests simulate a clean build, and then a fully cached build (test cases 1 and 2).
  * **master prompt:**

> following the red green refactor tdd method, create package `internal/scheduler`. write tests in `scheduler_test.go` to simulate a full execution pipeline: first verifying a clean build where everything executes, and then passing the resulting cache into a second run verifying a fully cached build where nothing executes. then, implement `scheduler.go` with a `Run(targets []*models.Target, cachePath string) error` function. this function must: 1. load the cache, 2. build the graph, 3. check for cycles, 4. run topological sort, 5. iterate through the sorted list using a FIFO queue, checking the HashTable for existing hashes, 6. mark uncached targets as executed and add their hashes to both the HashTable and AVL tree, and 7. save the updated cache state. make the tests pass. keep the code clean and simple with simple inline comments. run the tests. if they pass, commit and push the code with the message "feat: implement main build execution and caching loop". use context7

### phase 5: the command line interface

**mini phase 5.1: main application wiring**

  * **purpose:** build the actual CLI application that the user interacts with.
  * **definition of done:** `cmd/resolver/main.go` parses arguments and calls the scheduler.
  * **verification:** running `go run cmd/resolver/main.go build` successfully processes a local JSON configuration file.
  * **master prompt:**

> create the entry point for the application in `cmd/resolver/main.go`. use the standard `flag` or `os.Args` library to set up a simple CLI. implement a `build` command that reads a `data/build.json` file (which contains a list of Target structs), reads `data/cache.json`, and passes them into the `scheduler.Run` function. implement an `init` command that generates dummy `build.json` and `cache.json` files in a `data` folder so the user has something to test with. add basic print statements to provide a user interface for the console. test it manually to ensure it compiles and runs. keep the code clean and simple with simple inline comments. once working, commit and push the code with the message "feat: implement command line interface and application entrypoint". use context7