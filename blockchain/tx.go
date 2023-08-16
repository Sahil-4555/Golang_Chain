package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/Sahil-4555/Golang_Chain/wallet"
)

// TxOutput represents an output in a transaction
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// TxOutputs holds a collection of transaction outputs
type TxOutputs struct {
	Outputs []TxOutput
}

// TxInput represents an input in a transaction
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

// Check if a transaction input uses a specific public key hash
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock a transaction output to a specific address
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// Check if a transaction output is locked with a specific key
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// Create a new transaction output with a specific value and address
func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// Serialize the transaction outputs
func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(outs)
	Handle(err)

	return buffer.Bytes()
}

// Deserialize transaction outputs from serialized data
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	Handle(err)

	return outputs
}
