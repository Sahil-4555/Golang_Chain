package blockchain

import (
	"bytes"         
	"crypto/sha256" 
	"encoding/binary" 
	"log" 
	"math" 
	"math/big" 
)

// Set the difficulty level for mining. Higher values make mining harder.
const Difficulty = 18

// Define a structure called "ProofOfWork" for managing the mining process.
type ProofOfWork struct {
	Block  *Block      // The block to be mined.
	Target *big.Int    // The target value to meet for mining.
}

// Create a new ProofOfWork instance with a given block.
func NewProof(b *Block) *ProofOfWork {
	// Calculate the target value based on the difficulty level.
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	// Initialize the ProofOfWork with the block and target.
	pow := &ProofOfWork{b, target}
	return pow
}

// Prepare the data for mining by combining block information and nonce.
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

// Perform the mining process to find a valid nonce and hash.
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	// Keep incrementing the nonce until a valid hash is found.
	for nonce < math.MaxInt32 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])

		// Check if the calculated hash is less than the target.
		if intHash.Cmp(pow.Target) == -1 {
			break // Found a valid nonce that meets the target.
		} else {
			nonce++ // Increment nonce and try again.
		}
	}
	return nonce, hash[:]
}

// Validate checks if a block's nonce satisfies the mining target.
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	// Check if the calculated hash is less than the target.
	return intHash.Cmp(pow.Target) == -1
}

// Convert an int64 to its byte representation.
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err) // Handle any errors by panicking.
	}
	return buff.Bytes()
}
