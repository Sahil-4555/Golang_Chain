package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Define the file path for storing wallet data
const walletFile = "./tmp/wallets_%s.data"

// Wallets represents a collection of wallets.
type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallets initializes and loads wallets from a file.
func CreateWallets(nodeId string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	// Load wallets from a file (if it exists)
	err := wallets.LoadFile(nodeId)

	return &wallets, err
}

// AddWallet creates a new wallet and adds it to the collection.
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	// Add the wallet to the collection
	ws.Wallets[address] = wallet

	return address
}

// GetAllAddresses returns a list of all wallet addresses.
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	// Iterate through wallets and collect their addresses
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet retrieves a specific wallet by its address.
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// LoadFile reads wallet data from a file and populates the collection.
func (ws *Wallets) LoadFile(nodeId string) error {
	// Generate the file path for the wallet data
	walletFile := fmt.Sprintf(walletFile, nodeId)

	// Check if the wallet file exists
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err // File does not exist
	}

	var wallets Wallets

	// Read the content of the wallet file
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	// Register the elliptic curve to decode the data
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	// Decode the wallet data into the 'wallets' variable
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	// Populate the current Wallets collection with the loaded data
	ws.Wallets = wallets.Wallets

	return nil
}

// SaveFile encodes and saves the wallet data to a file.
func (ws *Wallets) SaveFile(nodeId string) {
	var content bytes.Buffer

	// Generate the file path for the wallet data
	walletFile := fmt.Sprintf(walletFile, nodeId)

	// Register the elliptic curve to encode the data
	gob.Register(elliptic.P256())

	// Create an encoder and encode the wallet data
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	// Write the encoded data to the wallet file
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
