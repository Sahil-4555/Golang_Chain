package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Proof of Work (PoW) algorithm implementation

// Difficulty level for mining
const Difficulty = 18

// ProofOfWork represents a proof of work for a block
type ProofOfWork struct {
	// The block to work on
	Block  *Block
	// The target value for proof of work
	Target *big.Int
}

// NewProof creates a new ProofOfWork instance with the given block
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	// Left-shift target to create the desired difficulty level
	target.Lsh(target, uint(256-Difficulty))

	// Create a ProofOfWork instance
	pow := &ProofOfWork{b, target}

	// Return the created ProofOfWork instance
	return pow
}

// InitData prepares the data for proof of work calculation with a specific nonce
func (pow *ProofOfWork) InitData(nonce int) []byte {

	// Concatenate previous hash, data, nonce in hexadecimal format, difficulty in hexadecimal format
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	// Return the prepared data
	return data
}

// Run performs the proof of work calculation and returns the nonce and resulting hash
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	// Initialize nonce to 0
	nonce := 0

	// Iterate until the maximum int32 value
	for nonce < math.MaxInt32 {
		// Get the initialized data
		data := pow.InitData(nonce)
		// Calculate the SHA-256 hash of the data
		hash = sha256.Sum256(data)

		// Print the hash in hexadecimal format
		fmt.Printf("\r%x", hash)
		// Convert the hash to a big integer
		intHash.SetBytes(hash[:])
		// Check if the calculated hash is less than the target
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			// If not, increment the nonce and continue the loop
			nonce++
		}

	}
	fmt.Println()

	// Return the final nonce and hash
	return nonce, hash[:]
}

// Validate checks if the proof of work is valid
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	// Get the initialized data
	data := pow.InitData(pow.Block.Nonce)

	// Calculate the SHA-256 hash of the data
	hash := sha256.Sum256(data)

	// Convert the hash to a big integer
	intHash.SetBytes(hash[:])

	// Check if the hash is less than the target
	return intHash.Cmp(pow.Target) == -1
}

// ToHex converts a number to a byte slice in hexadecimal format
func ToHex(num int64) []byte {

	// Create a new buffer for writing binary data
	buff := new(bytes.Buffer)

	// Write the number in binary to the buffer
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}
	// Return the bytes in the buffer
	return buff.Bytes()
}