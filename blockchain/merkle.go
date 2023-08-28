package blockchain

import "crypto/sha256"

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleNode creates a new Merkle node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	// Calculate hash based on left, right, and data
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	// Assign left and right children
	node.Left = left
	node.Right = right

	return &node
}

// NewMerkleTree creates a new Merkle tree from a list of data
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	// Ensure even number of elements in data
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	// Create initial nodes from data
	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, *node)
	}

	// Build the Merkle tree levels
	for i := 0; i < len(data)/2; i++ {
		var level []MerkleNode

		// Combine nodes to create parent nodes for the next level
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	// Create the Merkle tree with the root node
	tree := MerkleTree{&nodes[0]}

	return &tree
}
