package cache

// HashTable stores artifact hashes by target ID.
type HashTable struct {
	items map[string]string
}

// NewHashTable creates an empty hash table.
func NewHashTable() *HashTable {
	return &HashTable{
		items: make(map[string]string),
	}
}

// Set stores a hash for the given target ID.
func (h *HashTable) Set(id, hash string) {
	h.items[id] = hash
}

// Get retrieves a hash for the given target ID.
func (h *HashTable) Get(id string) (string, bool) {
	hash, ok := h.items[id]
	return hash, ok
}
