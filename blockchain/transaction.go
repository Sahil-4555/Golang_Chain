package blockchain

import (
	"bytes"         
	"crypto/sha256" 
	"encoding/gob" 
	"encoding/hex" 
	"fmt" 
	"log" 
)

// Transaction represents a transaction within the blockchain
type Transaction struct {
	ID      []byte     // The unique ID of the transaction
	Inputs  []TxInput  // The list of input transactions
	Outputs []TxOutput // The list of output transactions
}

// TxOutput represents an output of a transaction
type TxOutput struct {
	Value  int    // The value of the transaction output
	PubKey string // The public key associated with the output
}

// TxInput represents an input to a transaction
type TxInput struct {
	ID  []byte // The ID of the transaction input
	Out int    // The output index of the input
	Sig string // The signature of the input
}

// SetID generates a unique ID for the transaction based on its content
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// CoinbaseTx creates a special coinbase transaction to reward miners
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

// NewTransaction creates a new transaction between two addresses
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

// IsCoinbase checks if a transaction is a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// CanUnlock checks if a transaction input can be unlocked using the provided data
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// CanBeUnlocked checks if a transaction output can be unlocked using the provided data
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
