package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Block represents a block in the blockchain.
type Block struct {
	// Hash of the block
	Hash     []byte
	// Data stored in the block
	Data     []byte
	// Hash of the previous block
	PrevHash []byte
	// Nonce value used for proof of work
	Nonce    int
}

// CreateBlock creates a new block with the provided data and previous hash.
func CreateBlock(data string, prevHash []byte) *Block {
	// Initialize a new block
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	// Create a proof of work instance
	pow := NewProof(block)
	// Calculate the nonce and hash using proof of work
	nonce, hash := pow.Run()

	// Set the block's hash
	block.Hash = hash[:]
	// Set the block's nonce
	block.Nonce = nonce

	// Return the created block
	return block
}

// Genesis returns the genesis block of the blockchain
func Genesis() *Block {
	// Create the genesis block with empty data and hash
	return CreateBlock("Genesis", []byte{})
}

// Serialize encodes a block into a byte slice.
func (b *Block) Serialize() []byte {
	// Create a buffer to store the serialized data
	var res bytes.Buffer
	// Initialize a new encoder with the buffer
	encoder := gob.NewEncoder(&res)

	// Encode the block into the buffer
	err := encoder.Encode(b)

	// Handle any encoding errors
	Handle(err)

	// Return the serialized data as a byte slice
	return res.Bytes()
}

// Deserialize decodes a byte slice into a block.
func Deserialize(data []byte) *Block {
	// Create a new block instance
	var block Block

	// Initialize a new decoder with the byte slice
	decoder := gob.NewDecoder(bytes.NewReader(data))

	// Decode the data into the block instance
	err := decoder.Decode(&block)

	// Handle any decoding errors
	Handle(err)

	// Return the decoded block instance
	return &block
}

// Handle logs and panics if an error is provided.
func Handle(err error) {
	if err != nil {
		// Log the error and panic to halt program execution
		log.Panic(err)
	}
}