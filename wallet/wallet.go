package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4    // Length of the checksum in bytes
	version        = byte(0x00) // Version byte for address
)

// Wallet represents a cryptocurrency wallet containing private and public keys.
type Wallet struct {
	PrivateKey ecdsa.PrivateKey // Private key of the wallet
	PublicKey  []byte           // Public key of the wallet
}

// Address generates and returns the address associated with the wallet.
func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey) // Calculate the public key hash
	versionedHash := append([]byte{version}, pubHash...) // Add version to the hash
	checksum := Checksum(versionedHash) // Calculate the checksum
	fullHash := append(versionedHash, checksum...) // Combine version, hash, and checksum
	address := Base58Encode(fullHash) // Encode the full hash using Base58 encoding

	return address // Return the generated address
}

// NewKeyPair generates a new private and public key pair using elliptic curve cryptography.
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256() // Use elliptic curve P-256

	private, err := ecdsa.GenerateKey(curve, rand.Reader) // Generate a new private key
	if err != nil {
		log.Panic(err) // If there's an error, panic with the error message
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...) // Concatenate X and Y components of the public key
	return *private, pub // Return the private key and concatenated public key
}

// MakeWallet creates a new wallet by generating a private and public key pair.
func MakeWallet() *Wallet {
	private, public := NewKeyPair() // Generate a new key pair
	wallet := Wallet{private, public} // Create a new wallet using the keys

	return &wallet // Return a pointer to the created wallet
}

// PublicKeyHash calculates the hash of the public key using SHA-256 and RIPEMD-160.
func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey) // Hash the public key using SHA-256

	hasher := ripemd160.New() // Create a new RIPEMD-160 hasher
	_, err := hasher.Write(pubHash[:]) // Write the hashed public key to the hasher
	if err != nil {
		log.Panic(err) // If there's an error, panic with the error message
	}

	publicRipMD := hasher.Sum(nil) // Get the final RIPEMD-160 hash

	return publicRipMD // Return the RIPEMD-160 hash
}

// Checksum calculates a checksum for a given payload using double SHA-256 hashing.
func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload) // Calculate the first hash using SHA-256
	secondHash := sha256.Sum256(firstHash[:]) // Calculate the second hash using SHA-256

	return secondHash[:checksumLength] // Return the first `checksumLength` bytes of the second hash
}
