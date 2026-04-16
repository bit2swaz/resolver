package cache

import (
	"fmt"
	"testing"

	"github.com/bit2swaz/resolver/internal/models"
)

var hashSink string
var nodeSink *models.AVLNode

func BenchmarkHashTableLookup(b *testing.B) {
	for _, size := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			keys := benchmarkKeys(size)
			table := NewHashTable()
			for _, key := range keys {
				table.Set(key, "hash-"+key)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := keys[i%len(keys)]
				hash, ok := table.Get(key)
				if !ok {
					b.Fatalf("expected key %s to exist", key)
				}

				hashSink = hash
			}
		})
	}
}

func BenchmarkAVLLookup(b *testing.B) {
	for _, size := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			keys := benchmarkKeys(size)
			var root *models.AVLNode
			for _, key := range keys {
				root = Insert(root, key, "hash-"+key)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := keys[i%len(keys)]
				node := Search(root, key)
				if node == nil {
					b.Fatalf("expected key %s to exist", key)
				}

				nodeSink = node
			}
		})
	}
}

func BenchmarkHashTableBuild(b *testing.B) {
	for _, size := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			keys := benchmarkKeys(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				table := NewHashTable()
				for _, key := range keys {
					table.Set(key, "hash-"+key)
				}

				hashSink, _ = table.Get(keys[len(keys)-1])
			}
		})
	}
}

func BenchmarkAVLBuild(b *testing.B) {
	for _, size := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			keys := benchmarkKeys(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var root *models.AVLNode
				for _, key := range keys {
					root = Insert(root, key, "hash-"+key)
				}

				nodeSink = Search(root, keys[len(keys)-1])
			}
		})
	}
}

func benchmarkKeys(size int) []string {
	keys := make([]string, size)
	for index := range keys {
		keys[index] = fmt.Sprintf("target-%06d", index)
	}

	return keys
}
