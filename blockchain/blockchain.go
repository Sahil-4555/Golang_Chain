package blockchain

import (
	// Import required packages and libraries
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

// Define constants for database paths and the genesis data
const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

// Define the BlockChain structure to hold blockchain data
type BlockChain struct {
	LastHash []byte     // Hash of the last block in the blockchain
	Database *badger.DB // Database instance for block storage
}

// Define the BlockChainIterator structure for iterating through blocks
type BlockChainIterator struct {
	CurrentHash []byte     // Current hash being iterated
	Database    *badger.DB // Database instance for block retrieval
}

// DBexists checks if the database file exists
func DBexists() bool {
	// Check if the database file exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// ContinueBlockChain initializes an existing blockchain instance
func ContinueBlockChain(address string) *BlockChain {
	// Check if a blockchain already exists, if not, exit
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	// Initialize variables to store the last hash and database options
	var lastHash []byte
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Truncate = true // Truncate the value log to start fresh.
	opts.ValueLogLoadingMode = options.FileIO // Use FileIO mode for value log.

	// Open the Badger database with specified options
	db, err := badger.Open(opts)
	Handle(err)

	// Fetch the last hash from the database
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})
	Handle(err)

	// Create and return a BlockChain instance
	chain := BlockChain{lastHash, db}
	return &chain
}

// InitBlockChain creates a new blockchain instance and genesis block
func InitBlockChain(address string) *BlockChain {
	// Initialize variables to store the last hash and database options
	var lastHash []byte

	// If a blockchain already exists, print a message and exit
	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	// Initialize Badger database options
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Truncate = true // Truncate the value log to start fresh.
	opts.ValueLogLoadingMode = options.FileIO // Use FileIO mode for value log.

	// Open the Badger database with specified options
	db, err := badger.Open(opts)
	Handle(err)

	// Create a Coinbase transaction for the genesis block
	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		// Store the genesis block's hash and last hash in the database
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	// Create and return a BlockChain instance
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// AddBlock appends a new block with transactions to the blockchain
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	// Verify each transaction in the list
	var lastHash []byte
	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) != true {
			log.Panic("Invalid Transaction")
		}
	}

	// Fetch the last hash from the database
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})
	Handle(err)

	// Create a new block using the transactions and last hash
	newBlock := CreateBlock(transactions, lastHash)

	// Update the database with the new block and its hash
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

// Iterator creates and returns an iterator to traverse the blockchain
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

// Next advances the iterator and returns the next block
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// Fetch the current block's encoded data from the database
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	// Update the current hash to the previous block's hash
	iter.CurrentHash = block.PrevHash

	return block
}

// FindUnspentTransactions finds unspent transactions for a given public key hash
func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	// Initialize variables to hold unspent transactions and spent transaction outputs
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)

	// Create an iterator to traverse the blockchain
	iter := chain.Iterator()

	// Loop through the blockchain and transactions to find unspent transactions
	for {
		block := iter.Next()

		// Loop through transactions in the block
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// Loop through transaction outputs to find unspent ones
		Outputs:
			for outIdx, out := range tx.Outputs {
				// Check if the output is marked as spent
				if spentTXOs[txID] != nil {
					// Check if the output index is in the spent outputs list
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// Check if the output is locked with the provided public key hash
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			// Check inputs of non-coinbase transactions to update spent outputs
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		// Stop the loop if we reach the genesis block
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

// FindUTXO finds unspent transaction outputs for a given public key hash
func (chain *BlockChain) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput

	// Find unspent transactions for the given public key hash
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)

	// Loop through unspent transactions and their outputs to find UTXOs
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// FindSpendableOutputs finds spendable outputs for a given amount
func (chain *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

	// Loop through unspent transactions to find spendable outputs
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		// Loop through transaction outputs to find spendable ones
		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				// Stop when accumulated amount is sufficient
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

// FindTransaction searches for a transaction by its ID
func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	// Traverse the blockchain to find the desired transaction
	for {
		block := iter.Next()

		// Loop through transactions in the block
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		// Stop the loop if we reach the genesis block
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

// SignTransaction signs a transaction using a private key
func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	// Loop through transaction inputs to find previous transactions
	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	// Sign the transaction using the private key and previous transactions
	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies the validity of a transaction
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	// Loop through transaction inputs to find previous transactions
	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	// Verify the transaction using previous transactions
	return tx.Verify(prevTXs)
}
