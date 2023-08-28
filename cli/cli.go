package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/Sahil-4555/Golang_Chain/blockchain"
	"github.com/Sahil-4555/Golang_Chain/wallet"
)

// CommandLine represents the command line interface
type CommandLine struct{}

// printUsage prints the usage instructions for the CLI
func (cli *CommandLine) printUsage() {
	// Display commands and their descriptions
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
	fmt.Println(" createwallet - Creates a new Wallet")
	fmt.Println(" listaddresses - Lists the addresses in our wallet file")
	fmt.Println(" reindexutxo - Rebuilds the UTXO set")
}

// validateArgs validates the command line arguments
func (cli *CommandLine) validateArgs() {
	// Check if a command is provided
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// reindexUTXO reindexes the UTXO set
func (cli *CommandLine) reindexUTXO() {
	// Load the blockchain and UTXO set
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}

	// Reindex the UTXO set
	UTXOSet.Reindex()

	// Count and display the number of transactions in the UTXO set
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

// listAddresses lists the addresses in the wallet
func (cli *CommandLine) listAddresses() {
	// Load wallets and retrieve all addresses
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	// Print each address
	for _, address := range addresses {
		fmt.Println(address)
	}
}

// createWallet creates a new wallet
func (cli *CommandLine) createWallet() {
	// Load wallets, add a new wallet, and save the file
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	// Print the newly created address
	fmt.Printf("New address is: %s\n", address)
}

// printChain prints the blocks in the blockchain
func (cli *CommandLine) printChain() {
	// Load the blockchain and iterate through blocks
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		// Print block information
		block := iter.Next()
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		
		// Validate Proof of Work
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

		// Print each transaction in the block
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		// Exit loop when at the genesis block
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// createBlockChain creates a new blockchain
func (cli *CommandLine) createBlockChain(address string) {
	// Validate the provided address
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")	
	}

	// Initialize the blockchain with the provided address
	chain := blockchain.InitBlockChain(address)
	defer chain.Database.Close()

	// Reindex the UTXO set
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	fmt.Println("Finished!")
}

// getBalance gets the balance of an address
func (cli *CommandLine) getBalance(address string) {
	// Validate the provided address
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")	
	}

	// Load the blockchain and UTXO set, then calculate the balance
	chain := blockchain.ContinueBlockChain(address)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// send sends coins from one address to another
func (cli *CommandLine) send(from, to string, amount int) {
	// Validate source and destination addresses
	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not Valid")	
	}
	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not Valid")	
	}

	// Load the blockchain and UTXO set
	chain := blockchain.ContinueBlockChain(from)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	// Create transaction and coinbase transaction, then add to a block
	tx := blockchain.NewTransaction(from, to, amount, &UTXOSet)
	cbTx := blockchain.CoinbaseTx(from, "")
	block := chain.AddBlock([]*blockchain.Transaction{cbTx, tx})

	// Update the UTXO set
	UTXOSet.Update(block)
	fmt.Println("Success!")
}

// Run executes the command line interface
func (cli *CommandLine) Run() {
	// Validate the provided command
	cli.validateArgs()

	// Define flag sets for each command
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	// Define flag variables for each command
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	// Parse the command line arguments
	switch os.Args[1] {
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	// Execute the appropriate command based on parsed input
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
