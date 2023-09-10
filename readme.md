# Golang_Chain

Golang_Chain is a simple blockchain implementation in Go (Golang). It serves as a foundation for building decentralized applications, exploring blockchain technology, and understanding the core concepts of blockchain networks.

## Functionality

The Golang_Chain project includes the following core functionality:

1. **Blockchain:** Implements a basic blockchain structure for storing transactions and blocks.
2. **Consensus Algorithm:** Uses a simple Proof of Work (PoW) consensus mechanism for block validation.
3. **Transactions:** Supports the creation, validation, and storage of transactions within the blockchain.
4. **Wallets:** Provides functionality to create, manage, and store cryptocurrency wallets.
5. **Networking:** Enables peer-to-peer communication for nodes in the network.

## Key Features

- **Security:** Transactions are secured through digital signatures, and the blockchain is tamper-resistant.
- **Decentralization:** Multiple nodes can participate in the network, making it decentralized and resilient.
- **Transaction Verification:** Ensures that transactions are valid and consistent with the blockchain's state.
- **Mining:** Miners can add new blocks to the blockchain, earning rewards for their efforts.
- **Wallet Management:** Users can create and manage wallets for storing cryptocurrencies.

## Example Usage

```
// Example: Create a new blockchain
chain := InitBlockChain("your_address", "your_node_id")

// Example: Create a money transaction
tx := NewTransaction("sender_address", "recipient_address", amount, chain)

// Example: Add the transaction to a new block and mine it
newBlock := chain.MineBlock([]*Transaction{tx})
```

## Getting Started

1. Clone this repository.
2. Install Go and BadgerDB.
3. Run `go run main.go` to start a blockchain node.

