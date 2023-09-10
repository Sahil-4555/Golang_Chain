package blockchain

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/dgraph-io/badger"
)

// Define constants for UTXO prefix and prefix length.
var (
	utxoPrefix   = []byte("utxo-") // Prefix for Unspent Transaction Outputs (UTXO).
	prefixLength = len(utxoPrefix)  // Length of the UTXO prefix.
)

// UTXOSet represents the Unspent Transaction Outputs set and its associated blockchain.
type UTXOSet struct {
	Blockchain *BlockChain // The blockchain to which this UTXO set belongs.
}

// FindSpendableOutputs finds and returns unspent transaction outputs that can be spent to reach the desired amount.
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int) // Create a map to store spendable outputs.
	accumulated := 0                     // Initialize the accumulated amount to zero.
	db := u.Blockchain.Database           // Get the BadgerDB database associated with the blockchain.

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions // Create iterator options.
		it := txn.NewIterator(opts)           // Create a new iterator.
		defer it.Close()                      // Close the iterator when done.

		// Iterate through UTXOs in the database.
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item()       // Get the current item (UTXO).
			k := item.Key()         // Get the key (transaction ID).
			v, err := item.Value()  // Get the value (serialized outputs).
			Handle(err)             // Handle any errors.
			k = bytes.TrimPrefix(k, utxoPrefix) // Remove the UTXO prefix to get the transaction ID.
			txID := hex.EncodeToString(k)        // Convert the transaction ID to hexadecimal.
			outs := DeserializeOutputs(v)        // Deserialize the UTXO outputs.

			// Iterate through the outputs to find spendable ones.
			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value              // Increase the accumulated amount.
					unspentOuts[txID] = append(unspentOuts[txID], outIdx) // Store the spendable output.
				}
			}
		}
		return nil
	})
	Handle(err) // Handle any errors.

	return accumulated, unspentOuts // Return the accumulated amount and spendable outputs.
}

// FindUnspentTransactions finds and returns unspent transaction outputs associated with a given public key hash.
func (u UTXOSet) FindUnspentTransactions(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput // Create a slice to store unspent transaction outputs.

	db := u.Blockchain.Database // Get the BadgerDB database associated with the blockchain.

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions // Create iterator options.
		it := txn.NewIterator(opts)           // Create a new iterator.
		defer it.Close()                      // Close the iterator when done.

		// Iterate through UTXOs in the database.
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item()      // Get the current item (UTXO).
			v, err := item.Value() // Get the value (serialized outputs).
			Handle(err)            // Handle any errors.
			outs := DeserializeOutputs(v) // Deserialize the UTXO outputs.

			// Iterate through the outputs to find unspent ones associated with the provided public key hash.
			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out) // Append the unspent output to the slice.
				}
			}
		}
		return nil
	})
	Handle(err) // Handle any errors.

	return UTXOs // Return the unspent transaction outputs.
}

// CountTransactions counts the number of transactions in the UTXO set.
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.Database // Get the BadgerDB database associated with the blockchain.
	counter := 0                // Initialize the transaction counter to zero.

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions // Create iterator options.
		it := txn.NewIterator(opts)           // Create a new iterator.
		defer it.Close()                      // Close the iterator when done.

		// Iterate through UTXOs in the database and increment the counter.
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			counter++
		}
		return nil
	})

	Handle(err) // Handle any errors.

	return counter // Return the number of transactions in the UTXO set.
}

// Reindex rebuilds the UTXO set by deleting the existing UTXOs and adding new ones from the blockchain.
func (u UTXOSet) Reindex() {
	db := u.Blockchain.Database // Get the BadgerDB database associated with the blockchain.

	// Delete existing UTXOs with the specified prefix.
	u.DeleteByPrefix(utxoPrefix)

	// Find the UTXO set from the blockchain.
	UTXO := u.Blockchain.FindUTXO()

	err := db.Update(func(txn *badger.Txn) error {
		// Iterate through the UTXOs and store them in the database with the appropriate key.
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			Handle(err)
			key = append(utxoPrefix, key...)

			err = txn.Set(key, outs.Serialize())
			Handle(err)
		}

		return nil
	})
	Handle(err) // Handle any errors.
}

// Update updates the UTXO set based on a new block by removing spent outputs and adding new ones.
func (u *UTXOSet) Update(block *Block) {
	db := u.Blockchain.Database // Get the BadgerDB database associated with the blockchain.

	err := db.Update(func(txn *badger.Txn) error {
		// Iterate through the transactions in the block.
		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false { // Skip coinbase transactions.
				for _, in := range tx.Inputs {
					updatedOuts := TxOutputs{}      // Create updated output set.
					inID := append(utxoPrefix, in.ID...) // Create input ID with prefix.
					item, err := txn.Get(inID)     // Get the item (output) from the database.
					Handle(err)                     // Handle any errors.
					v, err := item.Value()          // Get the value (serialized outputs).
					Handle(err)                     // Handle any errors.

					outs := DeserializeOutputs(v) // Deserialize the UTXO outputs.

					// Iterate through the outputs and exclude the spent one.
					for outIdx, out := range outs.Outputs {
						if outIdx != in.Out {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						if err := txn.Delete(inID); err != nil { // If no outputs remain, delete the entry.
							log.Panic(err)
						}
					} else {
						if err := txn.Set(inID, updatedOuts.Serialize()); err != nil {
							log.Panic(err)
						}
					}
				}
			}

			newOutputs := TxOutputs{} // Create a new output set for the transaction.
			for _, out := range tx.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			txID := append(utxoPrefix, tx.ID...) // Create transaction ID with prefix.
			if err := txn.Set(txID, newOutputs.Serialize()); err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	Handle(err) // Handle any errors.
}

// DeleteByPrefix deletes entries in the database with a specified prefix.
func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	// Define a function to delete keys in the database.
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := u.Blockchain.Database.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000 // Define the maximum number of keys to collect for deletion.
	u.Blockchain.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize) // Create a slice to collect keys for deletion.
		keysCollected := 0                               // Initialize the collected key counter to zero.
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)     // Copy the current key.
			keysForDelete = append(keysForDelete, key) // Append the key to the slice.
			keysCollected++                   // Increment the collected key counter.
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					log.Panic(err)
				}
				keysForDelete = make([][]byte, 0, collectSize) // Reset the keys slice.
				keysCollected = 0                               // Reset the key counter.
			}
		}
		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}
