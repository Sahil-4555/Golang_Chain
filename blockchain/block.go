package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
}

// HashTransactions calculates the hash of transactions in the block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

// CreateBlock creates a new block with transactions, previous hash, and a height
func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	// Create a new block with a timestamp, empty hash, provided transactions, previous hash, nonce 0, and height
	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	
	// Create a proof-of-work instance for this block
	pow := NewProof(block)
	
	// Run the proof-of-work algorithm to find a valid nonce and hash
	nonce, hash := pow.Run()

	// Set the calculated hash and nonce to the block
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis creates the first block (genesis block) with a coinbase transaction
func Genesis(coinbase *Transaction) *Block {
	// Create the genesis block with only the coinbase transaction, no previous hash, and height 0
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// Serialize converts a block into a byte slice
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	// Encode the block into bytes
	err := encoder.Encode(b)

	// Handle any errors
	Handle(err)

	return res.Bytes()
}

// Deserialize converts a byte slice into a block
func Deserialize(data []byte) *Block {
	var block Block

	// Create a decoder to read the data from bytes
	decoder := gob.NewDecoder(bytes.NewReader(data))

	// Decode the bytes into a block
	err := decoder.Decode(&block)

	// Handle any errors
	Handle(err)

	return &block
}

// Handle handles errors by logging and panicking
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
