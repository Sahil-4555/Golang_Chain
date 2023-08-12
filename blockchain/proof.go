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

// Difficulty level for mining
const Difficulty = 18

// ProofOfWork represents the proof of work algorithm for mining
type ProofOfWork struct {
	Block  *Block      // The block to be mined
	Target *big.Int    // The target value to meet for mining
}

// NewProof creates a new ProofOfWork instance
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty)) // Calculate the target value

	pow := &ProofOfWork{b, target}

	return pow
}

// InitData prepares the data for mining by combining block information and nonce
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

// Run performs the mining process to find a valid nonce and hash
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt32 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break // Found a valid nonce that meets the target
		} else {
			nonce++ // Increment nonce and try again
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// Validate checks if a block's nonce satisfies the mining target
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1 // Check if the hash is less than the target
}

// ToHex converts an int64 to its byte representation
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err) // Handle any errors by panicking

	}

	return buff.Bytes()
}
