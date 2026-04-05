package cache

import "github.com/bit2swaz/resolver/internal/models"

// Insert adds or updates a node in the AVL tree and keeps it balanced.
func Insert(root *models.AVLNode, key, hash string) *models.AVLNode {
	if root == nil {
		return &models.AVLNode{
			Key:    key,
			Hash:   hash,
			Height: 1,
		}
	}

	if key < root.Key {
		root.Left = Insert(root.Left, key, hash)
	} else if key > root.Key {
		root.Right = Insert(root.Right, key, hash)
	} else {
		root.Hash = hash
		return root
	}

	root.Height = 1 + max(height(root.Left), height(root.Right))
	balance := balanceFactor(root)

	if balance > 1 {
		if key > root.Left.Key {
			root.Left = rotateLeft(root.Left)
		}
		return rotateRight(root)
	}

	if balance < -1 {
		if key < root.Right.Key {
			root.Right = rotateRight(root.Right)
		}
		return rotateLeft(root)
	}

	return root
}

// Search looks up a node by key in the AVL tree.
func Search(root *models.AVLNode, key string) *models.AVLNode {
	if root == nil {
		return nil
	}

	if key < root.Key {
		return Search(root.Left, key)
	}

	if key > root.Key {
		return Search(root.Right, key)
	}

	return root
}

func height(node *models.AVLNode) int {
	if node == nil {
		return 0
	}

	return node.Height
}

func balanceFactor(node *models.AVLNode) int {
	if node == nil {
		return 0
	}

	return height(node.Left) - height(node.Right)
}

func rotateRight(node *models.AVLNode) *models.AVLNode {
	newRoot := node.Left
	transferredSubtree := newRoot.Right

	newRoot.Right = node
	node.Left = transferredSubtree

	node.Height = 1 + max(height(node.Left), height(node.Right))
	newRoot.Height = 1 + max(height(newRoot.Left), height(newRoot.Right))

	return newRoot
}

func rotateLeft(node *models.AVLNode) *models.AVLNode {
	newRoot := node.Right
	transferredSubtree := newRoot.Left

	newRoot.Left = node
	node.Right = transferredSubtree

	node.Height = 1 + max(height(node.Left), height(node.Right))
	newRoot.Height = 1 + max(height(newRoot.Left), height(newRoot.Right))

	return newRoot
}

func max(left, right int) int {
	if left > right {
		return left
	}

	return right
}