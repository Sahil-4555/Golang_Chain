package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/Sahil-4555/Golang_Chain/wallet"
)

// TxOutput represents an output of a transaction with a specific value and public key hash.
type TxOutput struct {
	Value      int    // The value (amount) of the output.
	PubKeyHash []byte // Hash of the public key.
}

// TxOutputs represents a collection of transaction outputs.
type TxOutputs struct {
	Outputs []TxOutput // Slice of TxOutput objects.
}

// TxInput represents an input to a transaction, including its ID, signature, and public key.
type TxInput struct {
	ID        []byte // Transaction ID.
	Out       int    // Output index within the referenced transaction.
	Signature []byte // Digital signature.
	PubKey    []byte // Public key.
}

// UsesKey checks if a TxInput uses the provided public key hash.
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock sets the PubKeyHash of a TxOutput based on an address.
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)       // Decode the address from Base58.
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]   // Remove version and checksum bytes.
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if a TxOutput is locked with the provided public key hash.
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput creates a new TxOutput with a specified value and locked to the given address.
func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address)) // Lock the output to the provided address.

	return txo
}

// Serialize encodes TxOutputs as a byte slice.
func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(outs)
	Handle(err)
	return buffer.Bytes()
}

// DeserializeOutputs decodes a byte slice into TxOutputs.
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	Handle(err)
	return outputs
}