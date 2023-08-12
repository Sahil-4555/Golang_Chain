package blockchain

import (
	"encoding/hex"       
	"fmt"                
	"os"                 
	"runtime"            
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"  
)

// Define constants for the database path, file, and genesis block data.
const (
	dbPath      = "./tmp/blocks"       // Directory path for database storage.
	dbFile      = "./tmp/blocks/MANIFEST"  // Path to the database file.
	genesisData = "First Transaction from Genesis"  // Initial block data.
)

// Define a struct named "BlockChain" to hold blockchain information.
type BlockChain struct {
	LastHash []byte          // The hash of the last block in the chain.
	Database *badger.DB      // The Badger database instance.
}

// Define a struct named "BlockChainIterator" to iterate over blocks.
type BlockChainIterator struct {
	CurrentHash []byte       // The hash of the current block.
	Database    *badger.DB   // The Badger database instance.
}

// Function to check if the database file exists.
func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// Function to continue an existing blockchain.
func ContinueBlockChain(address string) *BlockChain {
	// Check if the database file exists.
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit() // Exit the program if no blockchain exists.
	}

	var lastHash []byte

	// Configure Badger options for opening the database.
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Truncate = true // Truncate the value log to start fresh.
	opts.ValueLogLoadingMode = options.FileIO // Use FileIO mode for value log.

	// Open the Badger database.
	db, err := badger.Open(opts)
	Handle(err) // Handle any potential errors.

	// Retrieve the hash of the last block from the database.
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})
	Handle(err) // Handle any potential errors.

	// Create a new BlockChain instance with the retrieved hash and database.
	chain := BlockChain{lastHash, db}
	return &chain
}

// Function to initialize a new blockchain.
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	// Check if the blockchain already exists.
	if DBexists() {
		fmt.Println("Blockchain Already Exists")
		runtime.Goexit() // Exit the program if blockchain already exists.
	}

	// Configure Badger options for opening the database.
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Truncate = true // Truncate the value log to start fresh.
	opts.ValueLogLoadingMode = options.FileIO // Use FileIO mode for value log.

	// Open the Badger database.
	db, err := badger.Open(opts)
	Handle(err) // Handle any potential errors.

	// Create the genesis coinbase transaction.
	cbtx := CoinbaseTx(address, genesisData)
	// Create the genesis block with the coinbase transaction.
	genesis := Genesis(cbtx)

	// Update the database with the genesis block and its hash.
	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	Handle(err) // Handle any potential errors.

	// Create a new BlockChain instance with the genesis hash and database.
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// Function to add a new block with transactions to the blockchain.
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	// Retrieve the hash of the last block from the database.
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		return err
	})
	Handle(err) // Handle any potential errors.

	// Create a new block with the provided transactions.
	newBlock := CreateBlock(transactions, lastHash)

	// Update the database with the new block and its hash.
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err) // Handle any potential errors.
}

// Function to create an iterator over the blockchain.
func (chain *BlockChain) Iterator() *BlockChainIterator {
	// Create an iterator instance with the current hash and database.
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

// Function to retrieve the next block using the iterator.
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// Retrieve the next block's encoded data from the database.
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)
		return err
	})
	Handle(err) // Handle any potential errors.

	// Update the current hash for the next iteration.
	iter.CurrentHash = block.PrevHash
	return block
}

// Function to find unspent transactions for a specific address.
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int) // Track spent transaction outputs.

	// Create an iterator to navigate through the blockchain.
	iter := chain.Iterator()

	// Iterate over all blocks in the blockchain.
	for {
		block := iter.Next()

		// Iterate over transactions in the block.
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Iterate over transaction outputs.
		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					// Check if the output is already spent.
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs // Skip if output is spent.
						}
					}
				}
				// Check if output can be unlocked by the address.
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			// Check if the transaction is not a coinbase transaction.
			if tx.IsCoinbase() == false {
				// Iterate over transaction inputs.
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						// Update spent outputs with input references.
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		// Exit the loop when there are no more blocks.
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

// Function to find unspent transaction outputs (UTXOs) for an address.
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	// Find unspent transactions for the given address.
	unspentTransactions := chain.FindUnspentTransactions(address)

	// Iterate over unspent transactions and their outputs.
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// Function to find spendable outputs for a given amount and address.
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

	// Iterate over unspent transactions and their outputs.
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		// Iterate over transaction outputs.
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				// Update accumulated amount and spent outputs.
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				// Exit loop when the required amount is reached.
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	// Return the accumulated amount and the list of spendable outputs.
	return accumulated, unspentOuts
}
