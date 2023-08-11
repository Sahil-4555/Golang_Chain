package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger" 
	"github.com/dgraph-io/badger/options" 
)

const (
	dbPath = "./tmp/blocks" // Define the path for the database files
)

// BlockChain represents the blockchain.
type BlockChain struct {
	LastHash []byte      // Hash of the last block in the chain
	Database *badger.DB  // Badger database instance
}

// BlockChainIterator helps iterate through the blockchain.
type BlockChainIterator struct {
	CurrentHash []byte      // Current block's hash being iterated
	Database    *badger.DB  // Badger database instance
}

// InitBlockChain initializes and returns a new blockchain instance.
func InitBlockChain() *BlockChain {
	var lastHash []byte

	// Configure Badger database options
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Truncate = true // Truncate the value log
	opts.ValueLogLoadingMode = options.FileIO // Use FileIO mode for value log

	// Open Badger database with specified options
	db, err := badger.Open(opts)
	Handle(err) // Handle any errors

	// Update the database transaction to set up the genesis block if not present
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis() // Create the genesis block
			fmt.Println("Genesis block created")
			err = txn.Set(genesis.Hash, genesis.Serialize()) // Store the serialized genesis block
			Handle(err) // Handle any errors
			err = txn.Set([]byte("lh"), genesis.Hash) // Update the last hash key

			lastHash = genesis.Hash // Set the last hash to the genesis block's hash

			return err
		} else {
			item, err := txn.Get([]byte("lh")) // Retrieve the last hash from the database
			Handle(err) // Handle any errors
			lastHash, err = item.Value() // Get the value of the last hash
			return err // Return any errors
		}
	})

	Handle(err) // Handle any errors

	// Create a blockchain instance with the last hash and the database
	blockchain := BlockChain{lastHash, db}
	return &blockchain // Return a pointer to the blockchain instance
}

// AddBlock adds a new block to the blockchain.
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// Access the database in view mode to retrieve the last hash
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh")) // Retrieve the last hash from the database
		Handle(err) // Handle any errors
		lastHash, err = item.Value() // Get the value of the last hash

		return err // Return any errors
	})
	Handle(err) // Handle any errors

	// Create a new block with the provided data and last hash
	newBlock := CreateBlock(data, lastHash)

	// Update the database in a transaction to add the new block
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize()) // Store the serialized new block
		Handle(err) // Handle any errors
		err = txn.Set([]byte("lh"), newBlock.Hash) // Update the last hash key

		chain.LastHash = newBlock.Hash // Update the last hash in the blockchain

		return err // Return any errors
	})
	Handle(err) // Handle any errors
}

// Iterator returns a new iterator for the blockchain.
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter // Return the iterator
}

// Next returns the next block in the blockchain iterator.
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// Access the database in view mode to retrieve the current block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash) // Retrieve the block by its hash
		Handle(err) // Handle any errors
		encodedBlock, err := item.ValueCopy(nil) // Get the serialized block
		block = Deserialize(encodedBlock) // Deserialize the block from the value

		return err // Return any errors
	})
	Handle(err) // Handle any errors

	iter.CurrentHash = block.PrevHash // Update the current hash to the previous hash

	return block // Return the retrieved block
}
