package blockchain

// TxOutput represents an output of a transaction, indicating the amount and recipient.
type TxOutput struct {
	Value  int    // Amount of coins
	PubKey string // Public key or address of the recipient
}

// TxInput represents an input of a transaction, linking to a previous output (UTXO) and containing a signature.
type TxInput struct {
	ID  []byte // ID of the transaction that the input references
	Out int    // Index of the output within the referenced transaction
	Sig string // Signature to unlock the input (authorization)
}

// CanUnlock checks if an input can be unlocked using a provided data (signature).
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data // Compares the provided signature with the input's signature
}

// CanBeUnlocked checks if an output can be unlocked using a provided data (public key or address).
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data // Compares the provided data with the output's public key
}
