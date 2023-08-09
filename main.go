package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// BlockChain represents a sequence of blocks
type BlockChain struct {
	// Slice to store blocks in the blockchain
	blocks []*Block
}

// Block represents a single block in the blockchain
type Block struct {
	// Hash of the block (output of hashing algorithm)
	Hash     []byte	
	// Data stored in the block
	Data     []byte
	// Hash of the previous block
	PrevHash []byte
}

// DeriveHash calculates the hash of the block by hashing its data and previous hash
func (b *Block) DeriveHash() {
	// Concatenate data and previous hash
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	// Calculate SHA-256 hash
	hash := sha256.Sum256(info)
	// Assign the calculated hash to the block's Hash field
	b.Hash = hash[:]
}

// CreateBlock creates a new block with given data and previous hash
func CreateBlock(data string, prevHash []byte) *Block {
	// Initialize a new block with empty hash, provided data, and previous hash
	block := &Block{[]byte{}, []byte(data), prevHash}
	// Calculate and assign the hash for the new block
	block.DeriveHash()
	// Return the created block
	return block
}

// AddBlock appends a new block with given data to the blockchain
func (chain *BlockChain) AddBlock(data string) {
	// Get the most recent block in the chain
	prevBlock := chain.blocks[len(chain.blocks)-1]
	// Create a new block using provided data and hash of previous block
	new := CreateBlock(data, prevBlock.Hash)
	// Append the new block to the blockchain
	chain.blocks = append(chain.blocks, new)
}

// Genesis creates the first block with initial data and no previous hash
func Genesis() *Block {
	// Create the first block with predefined data and empty previous hash
	return CreateBlock("Starting Block", []byte{})
}

// InitBlockChain initializes a new blockchain with the genesis block
func InitBlockChain() *BlockChain {
	// Create a new blockchain with the genesis block
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	// Initialize a new blockchain
	chain := InitBlockChain()

	// Add the first block to the chain
	chain.AddBlock("First Block..")
	// Add the second block to the chain
	chain.AddBlock("Second Block..")
	// Add the third block to the chain
	chain.AddBlock("Third Block..")

	// Print information about each block in the blockchain
	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n\n", block.Hash)
	}
}