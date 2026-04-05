package cache

import (
	"testing"

	"github.com/bit2swaz/resolver/internal/models"
)

func TestInsertBalancesTreeWithLeftRotation(t *testing.T) {
	var root *models.AVLNode
	root = Insert(root, "a", "hash-a")
	root = Insert(root, "b", "hash-b")
	root = Insert(root, "c", "hash-c")

	if root == nil {
		t.Fatal("expected root node")
	}

	if root.Key != "b" {
		t.Fatalf("expected balanced root b, got %q", root.Key)
	}

	if root.Left == nil || root.Left.Key != "a" {
		t.Fatal("expected left child a after left rotation")
	}

	if root.Right == nil || root.Right.Key != "c" {
		t.Fatal("expected right child c after left rotation")
	}

	if root.Height != 2 {
		t.Fatalf("expected root height 2, got %d", root.Height)
	}
}

func TestInsertBalancesTreeWithRightRotation(t *testing.T) {
	var root *models.AVLNode
	root = Insert(root, "c", "hash-c")
	root = Insert(root, "b", "hash-b")
	root = Insert(root, "a", "hash-a")

	if root == nil {
		t.Fatal("expected root node")
	}

	if root.Key != "b" {
		t.Fatalf("expected balanced root b, got %q", root.Key)
	}

	if root.Left == nil || root.Left.Key != "a" {
		t.Fatal("expected left child a after right rotation")
	}

	if root.Right == nil || root.Right.Key != "c" {
		t.Fatal("expected right child c after right rotation")
	}

	if root.Height != 2 {
		t.Fatalf("expected root height 2, got %d", root.Height)
	}
}

func TestSearchReturnsInsertedNode(t *testing.T) {
	var root *models.AVLNode
	root = Insert(root, "m", "hash-m")
	root = Insert(root, "a", "hash-a")
	root = Insert(root, "z", "hash-z")
	root = Insert(root, "b", "hash-b")

	node := Search(root, "b")
	if node == nil {
		t.Fatal("expected to find node b")
	}

	if node.Hash != "hash-b" {
		t.Fatalf("expected hash-b, got %q", node.Hash)
	}

	missing := Search(root, "missing")
	if missing != nil {
		t.Fatal("expected missing node search to return nil")
	}
}