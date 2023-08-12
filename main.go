package main

import (
	"flag" 
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/Sahil-4555/Golang_Chain/blockchain" // Import the blockchain package
)

type CommandLine struct{} // Define a struct named CommandLine

func (cli *CommandLine) printUsage() { // Define a method to print usage instructions
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
}

func (cli *CommandLine) validateArgs() { // Define a method to validate command-line arguments
	if len(os.Args) < 2 {
		cli.printUsage() // Print usage instructions if not enough arguments
		runtime.Goexit() // Exit the program
	}
}

func (cli *CommandLine) printChain() { // Define a method to print information about the blockchain
	chain := blockchain.ContinueBlockChain("") // Continue the blockchain
	defer chain.Database.Close() // Close the database when done

	iter := chain.Iterator() // Get an iterator for the blockchain

	for {
		block := iter.Next() // Move to the next block

		fmt.Printf("Prev. hash: %x\n", block.PrevHash) // Print previous block's hash
		fmt.Printf("Hash: %x\n", block.Hash) // Print current block's hash
		pow := blockchain.NewProof(block) // Create a new proof-of-work instance
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate())) // Print if proof-of-work is valid
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break // Exit the loop if at the genesis block
		}
	}
}

func (cli *CommandLine) createBlockChain(address string) { // Define a method to create a blockchain
	chain := blockchain.InitBlockChain(address) // Initialize the blockchain with the provided address
	chain.Database.Close() // Close the database
	fmt.Println("Finished!") // Print a message when finished
}

func (cli *CommandLine) getBalance(address string) { // Define a method to get the balance of an address
	chain := blockchain.ContinueBlockChain(address) // Continue the blockchain with the provided address
	defer chain.Database.Close() // Close the database when done

	balance := 0 // Initialize the balance to zero
	UTXOs := chain.FindUTXO(address) // Get unspent transaction outputs for the address

	for _, out := range UTXOs {
		balance += out.Value // Sum up the values of unspent outputs
	}

	fmt.Printf("Balance of %s: %d\n", address, balance) // Print the balance
}

func (cli *CommandLine) send(from, to string, amount int) { // Define a method to send coins
	chain := blockchain.ContinueBlockChain(from) // Continue the blockchain with the sender's address
	defer chain.Database.Close() // Close the database when done

	tx := blockchain.NewTransaction(from, to, amount, chain) // Create a new transaction
	chain.AddBlock([]*blockchain.Transaction{tx}) // Add the transaction to a new block
	fmt.Println("Success!") // Print a success message
}

// Define a method to run the CLI
func (cli *CommandLine) run() { 
	// Validate command-line arguments
	cli.validateArgs() 

	// Create a new flag set for getting balance
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	// Create a flag set for creating a blockchain 
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	// Create a flag set for sending coins 
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError) 
	// Create a flag set for printing the blockchain
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError) 

	// Define a flag for address in getbalance command
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for") 
	// Define a flag for address in createblockchain command
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to") 
	// Define a flag for source wallet address in send command
	sendFrom := sendCmd.String("from", "", "Source wallet address") 
	// Define a flag for destination wallet address in send command
	sendTo := sendCmd.String("to", "", "Destination wallet address") 
	// Define a flag for sending amount in send command
	sendAmount := sendCmd.Int("amount", 0, "Amount to send") 

	switch os.Args[1] {
	case "getbalance":
		// Parse arguments for getbalance command
		err := getBalanceCmd.Parse(os.Args[2:]) 
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		// Parse arguments for createblockchain command
		err := createBlockchainCmd.Parse(os.Args[2:]) 
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		// Parse arguments for printchain command
		err := printChainCmd.Parse(os.Args[2:]) 
		if err != nil {
			log.Panic(err)
		}
	case "send":
		// Parse arguments for send command
		err := sendCmd.Parse(os.Args[2:]) 
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage() // Print usage instructions for invalid command
		runtime.Goexit() // Exit the program
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		// Call the method to get the balance
		cli.getBalance(*getBalanceAddress) 
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		// Call the method to create a blockchain
		cli.createBlockChain(*createBlockchainAddress) 
	}

	if printChainCmd.Parsed() {
		// Call the method to print the blockchain
		cli.printChain() 
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		// Call the method to send coins
		cli.send(*sendFrom, *sendTo, *sendAmount) 
	}
}

func main() {
	
	// Ensure the program exits with a status code 0
	defer os.Exit(0) 
	// Create an instance of CommandLine struct
	cli := CommandLine{} 
	// Call the method to run the CLI
	cli.run() 
}
