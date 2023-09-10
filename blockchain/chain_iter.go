package blockchain

// Import the "badger" package for working with the database.
import "github.com/dgraph-io/badger"

// Define a struct called "BlockChainIterator" that represents an iterator for the blockchain.
type BlockChainIterator struct {
	CurrentHash []byte     // Current hash being iterated.
	Database    *badger.DB // The database used for retrieval.
}

// Create a method for the "BlockChain" struct called "Iterator" that returns a new blockchain iterator.
func (chain *BlockChain) Iterator() *BlockChainIterator {
	// Initialize an iterator with the current hash and the database.
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	// Return the iterator.
	return iter
}

// Create a method for the "BlockChainIterator" struct called "Next" that retrieves the next block.
func (iter *BlockChainIterator) Next() *Block {
	var block *Block // Initialize a variable to hold the block.

	// Use a database view to retrieve data.
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash) // Get the data associated with the current hash.
		Handle(err)                             // Handle any potential errors.
		encodedBlock, err := item.Value()        // Retrieve the encoded block data.
		block = Deserialize(encodedBlock)       // Deserialize the block data into a block structure.

		return err
	})
	Handle(err) // Handle any potential errors.

	iter.CurrentHash = block.PrevHash // Update the current hash to the previous block's hash.

	return block // Return the retrieved block.
}
