package blockchain

import "crypto/sha256"

// Define a structure called "MerkleTree" for representing a Merkle tree.
type MerkleTree struct {
	RootNode *MerkleNode // The root node of the Merkle tree.
}

// Define a structure called "MerkleNode" for representing nodes in the Merkle tree.
type MerkleNode struct {
	Left  *MerkleNode // The left child node.
	Right *MerkleNode // The right child node.
	Data  []byte      // Data associated with the node.
}

// Create a new Merkle node with optional left and right child nodes and some data.
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	// Create a new Merkle node.
	node := MerkleNode{}

	// If there are no left and right children (leaf node), calculate the hash of the data.
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else { // If there are child nodes, concatenate their data and calculate the hash.
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	// Set the left and right child nodes.
	node.Left = left
	node.Right = right

	// Return the newly created node.
	return &node
}

// Create a new Merkle tree from a list of data items (represented as byte slices).
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode // Create an empty list of Merkle nodes.

	// If the number of data items is odd, duplicate the last item to make it even.
	if len(data) % 2 != 0 {
		data = append(data, data[len(data)-1])
	}

	// Create a Merkle node for each data item and add it to the list of nodes.
	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, *node)
	}

	// Repeatedly combine pairs of nodes until only the root node remains.
	for i := 0; i < len(data) / 2; i++ {
		var level []MerkleNode

		// Combine pairs of nodes into parent nodes.
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		nodes = level // Update the list of nodes for the next level of the tree.
	}

	// Create the Merkle tree with the root node.
	tree := MerkleTree{&nodes[0]}

	// Return the Merkle tree.
	return &tree
}
