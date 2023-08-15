package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/Sahil-4555/Golang_Chain/wallet"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	ID      []byte     // Unique identifier for the transaction
	Inputs  []TxInput  // List of transaction inputs
	Outputs []TxOutput // List of transaction outputs
}

// Hash calculates and returns the hash of the transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	// Create a copy of the transaction without the ID
	txCopy := *tx
	txCopy.ID = []byte{}

	// Calculate the SHA-256 hash of the serialized transaction
	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// Serialize serializes the transaction into a byte slice
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	// Encode the transaction into a byte buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// SetID sets the ID of the transaction by hashing its contents
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	// Encode the transaction into a byte buffer
	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	// Calculate the SHA-256 hash of the serialized transaction
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// CoinbaseTx creates a new coinbase transaction
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	// Create a transaction input and output for the coinbase transaction
	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(100, to)

	// Create and set the ID of the coinbase transaction
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.SetID()

	return &tx
}

// NewTransaction creates a new transaction from sender to receiver
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	// Retrieve sender's wallet and public key hash
	wallets, err := wallet.CreateWallets()
	Handle(err)
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	// Find spendable outputs to cover the desired amount
	acc, validOutputs := chain.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	// Create transaction inputs from valid spendable outputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Create transaction output for the receiver
	outputs = append(outputs, *NewTXOutput(amount, to))

	// Create transaction output for change if applicable
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	// Create the transaction, set its ID, and sign it
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	chain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

// IsCoinbase checks if a transaction is a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// Sign signs a transaction using a private key and previous transactions
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	// Validate previous transactions for each input
	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	// Create a trimmed copy of the transaction for signing
	txCopy := tx.TrimmedCopy()

	// Loop through inputs and sign each one
	for inID, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTX.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].PubKey = nil

		// Sign the transaction copy using the private key
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		// Set the signature for the input in the original transaction
		tx.Inputs[inID].Signature = signature
	}
}

// Verify verifies the validity of a transaction using previous transactions
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	// Validate previous transactions for each input
	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Previous transaction not correct")
		}
	}

	// Create a trimmed copy of the transaction for verification
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	// Loop through inputs and verify each one
	for inID, in := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTX.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].PubKey = nil

		// Extract components of the signature
		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		// Extract components of the public key
		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		// Reconstruct the public key
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		// Verify the signature using the reconstructed public key and transaction hash
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

// TrimmedCopy creates a trimmed copy of the transaction (without signatures)
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	// Create inputs for the trimmed copy
	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	// Create outputs for the trimmed copy
	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	// Create the trimmed copy of the transaction
	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// String returns a string representation of the transaction
func (tx Transaction) String() string {
	var lines []string

	// Add transaction ID to the representation
	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	// Loop through transaction inputs
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	// Loop through transaction outputs
	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	// Combine the lines into a single string
	return strings.Join(lines, "\n")
}
