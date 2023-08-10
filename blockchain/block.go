package blockchain

// BlockChain represents a sequence of blocks
type BlockChain struct {
	// A collection of blocks in the blockchain
	Blocks []*Block
}

// Block represents a single block in the blockchain
type Block struct {
	// Hash of the block's data and previous hash
	Hash     []byte
	// Data stored in the block
	Data     []byte
	// Hash of the previous block
	PrevHash []byte
	// A nonce value used for proof of work
	Nonce    int
}

// CreateBlock creates a new block with given data and previous hash
func CreateBlock(data string, prevHash []byte) *Block {

	// Initialize a new block with data, previous hash, and nonce set to 0
	block := &Block{[]byte{}, []byte(data), prevHash, 0}

	// Create a proof of work instance for the block
	pow := NewProof(block)

	// Perform proof of work computation
	nonce, hash := pow.Run()

	// Set the computed hash to the block's hash
	block.Hash = hash[:]

	// Set the computed nonce to the block's nonce
	block.Nonce = nonce

	// Return the created block with updated hash and nonce
	return block
}

// AddBlock appends a new block with given data to the blockchain
func (chain *BlockChain) AddBlock(data string) {

	// Get the most recent block in the chain
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	// Create a new block using provided data and hash of previous block
	new := CreateBlock(data, prevBlock.Hash)

	// Append the new block to the blockchain
	chain.Blocks = append(chain.Blocks, new)
}

// Genesis creates the first block with initial data and no previous hash
func Genesis() *Block {

	// Create the first block with predefined data and empty previous hash
	return CreateBlock("Starting...", []byte{})
}

// InitBlockChain initializes a new blockchain with the genesis block
func InitBlockChain() *BlockChain {

	// Create a new blockchain with the genesis block
	return &BlockChain{[]*Block{Genesis()}}
}