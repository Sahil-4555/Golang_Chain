package blockchain

import (
	"bytes"

	"github.com/Sahil-4555/Golang_Chain/wallet"
)

// TxOutput represents an output in a transaction
type TxOutput struct {
	Value      int    // Amount of cryptocurrency in the output
	PubKeyHash []byte // Hash of the public key used to lock the output
}

// TxInput represents an input in a transaction
type TxInput struct {
	ID        []byte // Transaction ID that this input references
	Out       int    // Index of the output within the referenced transaction
	Signature []byte // Signature that authorizes spending this input
	PubKey    []byte // Public key associated with the signature
}

// UsesKey checks if a transaction input uses a specific public key hash
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	// Calculate the public key hash of the input's public key
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	// Compare the calculated public key hash with the given public key hash
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock sets the PubKeyHash of a transaction output using an address
func (out *TxOutput) Lock(address []byte) {
	// Decode the Base58 address to get the public key hash
	pubKeyHash := wallet.Base58Decode(address)

	// Remove the version byte and checksum to get the actual public key hash
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	// Set the PubKeyHash of the output to the extracted public key hash
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if a transaction output is locked with a specific public key hash
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	// Compare the public key hash of the output with the given public key hash
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput creates a new transaction output with a specified value and address
func NewTXOutput(value int, address string) *TxOutput {
	// Create a new transaction output with the given value and nil PubKeyHash
	txo := &TxOutput{value, nil}

	// Lock the transaction output using the provided address
	txo.Lock([]byte(address))

	return txo
}
