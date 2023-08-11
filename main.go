package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/Sahil-4555/Golang_Chain/blockchain"
)

// CommandLine holds the blockchain and provides functions for command-line interaction.
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

// printUsage prints the usage instructions for command-line interactions.
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

// validateArgs checks if the required number of command-line arguments are provided.
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// addBlock adds a new block to the blockchain with the provided data.
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

// printChain iterates through the blockchain and prints information about each block.
func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// run is the central function for handling command-line arguments and actions.
func (cli *CommandLine) run() {
	cli.validateArgs()

	// Create a flag set for the "add" command
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// Create a flag set for the "print" command
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	// Define a flag to capture block data
	addBlockData := addBlockCmd.String("block", "", "Block data")

	// Parse the command-line arguments and execute corresponding actions
	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	// Perform actions based on the parsed commands
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

// The main function initializes the blockchain and starts the command-line interface.
func main() {
	// Ensure the program exits with status code 0
	defer os.Exit(0)
	// Initialize the blockchain
	chain := blockchain.InitBlockChain()
	// Close the database connection when the main function ends
	defer chain.Database.Close()

	// Create a CommandLine instance with the blockchain
	cli := CommandLine{chain}
	// Execute the command-line interface
	cli.run()
}