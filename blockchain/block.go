package blockchain 

import ( 
	"bytes"
	"crypto/sha256" 
	"encoding/gob" 
	"log" 
)

// Define a struct named "Block" to hold block information.
type Block struct {
	Hash         []byte          // Hash of the block.
	Transactions []*Transaction  // List of transactions in the block.
	PrevHash     []byte          // Hash of the previous block.
	Nonce        int             // A nonce value for proof-of-work.
}

// Function to calculate the hash of all transactions in the block.
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	// Collecting the IDs of each transaction and creating a combined hash.
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// Function to create a new block with transactions and previous hash.
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	// Create a block instance with initial values.
	block := &Block{[]byte{}, txs, prevHash, 0} 
	// Create a proof-of-work instance.
	pow := NewProof(block) 
	// Calculate nonce and hash using proof-of-work.
	nonce, hash := pow.Run() 
	// Set calculated hash.
	block.Hash = hash[:] 
	// Set calculated nonce.
	block.Nonce = nonce 
	// Return the new block instance.
	return block 
}

// Function to create the initial block (genesis block) with a coinbase transaction.
func Genesis(coinbase *Transaction) *Block {
	// Call CreateBlock with coinbase transaction.
	return CreateBlock([]*Transaction{coinbase}, []byte{}) 
}

// Function to serialize a block into a byte slice.
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	// Create a gob encoder.
	encoder := gob.NewEncoder(&res) 
	// Encode the block data.
	err := encoder.Encode(b) 
	// Handle errors and panic if needed.
	Handle(err) 
	// Return the serialized byte slice.
	return res.Bytes() 
}

// Function to deserialize a byte slice into a block.
func Deserialize(data []byte) *Block {
	var block Block
	// Create a gob decoder.
	decoder := gob.NewDecoder(bytes.NewReader(data)) 
	// Decode the data into a block instance.
	err := decoder.Decode(&block) 
	// Handle errors and panic if needed.
	Handle(err) 
	// Return the deserialized block instance.
	return &block 
}

// Function to handle errors and panic if needed.
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
