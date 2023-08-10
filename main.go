package main

import (
	"fmt"
	"strconv"

	"github.com/Sahil-4555/Golang_Chain/blockchain"
)

func main() {
	
	// Initialize a new blockchain
	chain := blockchain.InitBlockChain()

	// Add blocks to the blockchain
	chain.AddBlock("First Block...")
	chain.AddBlock("Second Block...")
	chain.AddBlock("Third Block...")

	// Loop through each block in the blockchain
	for _, block := range chain.Blocks {

		// Print the hash of the previous block
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)

		// Print the data stored in the block
		fmt.Printf("Data in Block: %s\n", block.Data)

		// Print the hash of the current block
		fmt.Printf("Hash: %x\n", block.Hash)

		// Create a new proof of work instance for the current block
		pow := blockchain.NewProof(block)

		// Validate the proof of work for the current block
		// Print the validation result as a boolean
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}
}