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

const walletFile = "./tmp/wallets.data" // Path to the wallet data file

// Wallets represents a collection of wallets.
type Wallets struct {
	Wallets map[string]*Wallet // Map of addresses to wallet pointers
}

// CreateWallets creates a new Wallets instance and loads wallets from the file if it exists.
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{} // Create a new empty Wallets instance
	wallets.Wallets = make(map[string]*Wallet) // Initialize the map

	err := wallets.LoadFile() // Load wallets from the file

	return &wallets, err // Return a pointer to the Wallets instance and any error
}

// AddWallet generates a new wallet, adds it to Wallets, and returns its address.
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet() // Create a new wallet
	address := fmt.Sprintf("%s", wallet.Address()) // Get the address of the wallet as a string

	ws.Wallets[address] = wallet // Add the wallet to the Wallets map

	return address // Return the wallet's address
}

// GetAllAddresses returns a list of all addresses present in Wallets.
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address) // Append each address to the addresses slice
	}

	return addresses // Return the list of addresses
}

// GetWallet retrieves a wallet using its address.
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address] // Return the wallet corresponding to the given address
}

// LoadFile reads wallet data from the file and populates the Wallets instance.
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err // Return an error if the file doesn't exist
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile) // Read the content of the wallet file
	if err != nil {
		return err // Return any error encountered during file reading
	}

	gob.Register(elliptic.P256()) // Register the elliptic curve type for decoding
	decoder := gob.NewDecoder(bytes.NewReader(fileContent)) // Create a decoder for decoding the file content
	err = decoder.Decode(&wallets) // Decode the content into the wallets instance
	if err != nil {
		return err // Return any error encountered during decoding
	}

	ws.Wallets = wallets.Wallets // Assign the loaded wallets to the Wallets instance

	return nil // Return nil to indicate success
}

// SaveFile encodes and saves the Wallets instance to the file.
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer // Create a buffer to hold the encoded data

	gob.Register(elliptic.P256()) // Register the elliptic curve type for encoding

	encoder := gob.NewEncoder(&content) // Create an encoder for encoding the Wallets instance
	err := encoder.Encode(ws) // Encode the Wallets instance into the buffer
	if err != nil {
		log.Panic(err) // Panic if there's an error during encoding
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644) // Write the buffer to the wallet data file
	if err != nil {
		log.Panic(err) // Panic if there's an error during file writing
	}
}
