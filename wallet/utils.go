package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

// Base58Encode encodes the input bytes using Base58 encoding.
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input) // Encode the input bytes using Base58 encoding

	return []byte(encode) // Return the encoded bytes
}

// Base58Decode decodes the input bytes using Base58 decoding.
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:])) // Decode the Base58 encoded input bytes
	if err != nil {
		log.Panic(err) // If there's an error during decoding, panic with the error message
	}

	return decode // Return the decoded bytes
}
