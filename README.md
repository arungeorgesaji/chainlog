# ChainLog

**ChainLog** is a decentralized, tamper-proof logging and audit system. It provides an immutable public record for any application that needs verifiable, timestamped data history where trust is distributed.

Instead of moving money, ChainLog moves **trust** - providing cryptographic proof that data existed at a certain time and hasn't been modified since, and it happens right in your terminal.

---

> **Development Status**: Peer-to-peer networking is currently under development. The current version operates as a single-node system with all core blockchain functionalities fully operational. However, it is not yet ready for real-world use until networking is implemented. Also although users can become validators by staking Logcoins, the actual network-wide validation of other users’ blocks is not yet functional.... 

## Table of Contents
1. [Core Concepts](#core-concepts)
2. [Getting Started](#getting-started)
3. [CLI Commands](#cli-commands)
4. [LogCoin Economy](#logcoin-economy)
5. [Architecture](#architecture)

---

## Core Concepts

### 1. BLOCKCHAIN BASICS
- **Blocks**: Data containers that get chained together
- **Mining**: Process to add new blocks to the chain  
- **Nodes**: Computers that run the network
- **Consensus**: All nodes agree on the same data

### 2. DATA STORAGE
- Store any information: messages, files, records
- Timestamped and ordered
- Immutable - can't be changed once added
- Transparent - everyone can see it

### 3. NETWORK
- **P2P**: Computers connect directly to each other
- **Broadcast**: New data gets sent to all nodes
- **Validation**: Nodes verify data is valid  
- **Sync**: All nodes stay updated

### 4. SECURITY
- **Cryptography**: Data is cryptographically signed
- **Proof-of-Work**: Mining prevents spam/attacks
- **Digital Signatures**: Prove who created data
- **Hashing**: Link blocks together securely

### 5. WALLET & IDENTITY
- **Key Pairs**: Public/private keys for each user
- **Addresses**: Unique IDs derived from public keys
- **Signing**: Users sign their transactions
- **Verification**: Anyone can verify signatures

### 6. TRANSACTIONS
- **Data Payload**: The actual information being stored
- **Sender**: Who created it
- **Signature**: Proof it's authentic
- **Timestamp**: When it was created

### 7. MINING
- **Proof-of-Work**: Solve math problem to add block
- **Difficulty**: Adjusts to keep block time consistent
- **Reward**: Motivation for miners (could be reputation)
- **Competition**: Miners race to find next block

---

## Getting Started

### Installation
```bash
# Clone the repository
git clone https://github.com/yourusername/chainlog.git
cd chainlog

# Build the project
go build -o chainlog

# Run ChainLog
./chainlog help
```

### Quick Start
```bash
# 1. Create a wallet
./chainlog wallet create

# 2. Start a node
./chainlog start 8080

# 3. Create and broadcast a transaction
./chainlog transaction create "My first log entry" 2
./chainlog transaction broadcast <tx_id>

# 4. Mine the block
./chainlog mine
```
---

## CLI Commands

### Node Management
```bash
start [port]                  # Start a node (default: 8080)
status                        # Show blockchain status
summary                       # Print full system summary
help                          # Show help message
```

### Wallet Operations
```bash
wallet create                 # Create a new wallet
wallet import <key>           # Import wallet from private key
wallet list                   # List all wallets
balance <address>             # Check account balance
```

### Transactions
```bash
transaction create <data> <fee>    # Create a transaction
transaction list                   # List pending transactions
transaction broadcast <tx_id>      # Broadcast transaction
transaction status <tx_id>         # Check transaction status
```

### Mining
```bash
mine                          # Mine pending transactions
difficulty check              # Show current vs. recommended difficulty
```

### Blockchain Operations
```bash
chain show                    # Display full blockchain
chain validate                # Validate blockchain integrity
save                          # Save blockchain and state to disk
load                          # Load blockchain and state from disk
```

### Economy & Staking
```bash
economy stats                 # Show LogCoin economics
fees                          # Show fee statistics
rewards                       # Show reward statistics
staking add <address> <amt>   # Stake LogCoins
staking list                  # List validators and stakes
```

## LogCoin Economy

ChainLog uses **LogCoins** as a utility token to secure the network and incentivize participation.

### How LogCoins Work

#### Mining Rewards
- Miners earn **10 LogCoins** for each valid block mined
- Rewards compensate miners for computational work
- Ensures blocks are produced consistently

#### Transaction Fees
- Users pay **1-5 LogCoins** to store data on-chain
- Fees prevent network spam and abuse
- Priority processing for higher fee transactions

#### Staking System
- Stake **100+ LogCoins** to become a validator
- Validators participate in consensus and governance
- Malicious behavior results in stake slashing

#### Governance
- LogCoin holders vote on protocol upgrades
- Voting power proportional to coin balance
- Community-driven development decisions

### Practical Benefits

**For Users:**
- Pay coins to store important data permanently
- Higher fees = faster confirmation times
- Build reputation through coin ownership

**For Miners:**
- Earn coins for securing the network
- Transaction fees provide additional income
- Compete to solve blocks and earn rewards

**For the Network:**
- Coins prevent Sybil attacks and spam
- Incentivizes decentralized participation
- Self-sustaining economic model

### Economic Model

#### Supply & Distribution
- Initial supply: **1,000,000 LogCoins** (genesis block)
- Block reward halving every **100,000 blocks** (~2 years)
- Maximum supply cap: **21,000,000 LogCoins**
- Circulating supply increases gradually through mining

#### Deflationary Mechanisms
- Transaction fees are partially burned (20%)
- Reduces total supply over time
- Creates scarcity as network usage grows

#### Use Cases Beyond Fees
- Data storage pricing in LogCoins
- Premium features (priority timestamps, batch uploads)
- Network service subscriptions
- Third-party integrations and APIs

This economic structure ensures ChainLog remains secure, decentralized, and sustainable while aligning incentives across all network participants.

---

## Architecture

### System Components
- **Blockchain Core**: Block validation, chain management, consensus
- **P2P Network**: Node discovery, peer communication, data propagation
- **Wallet System**: Key management, transaction signing, address generation
- **Mining Engine**: Proof-of-Work implementation, difficulty adjustment
- **Transaction Pool**: Mempool for pending transactions
- **Persistence Layer**: Blockchain and state storage

### Security Features
- SHA-256 cryptographic hashing
- ECDSA digital signatures
- Proof-of-Work consensus
- Chain validation and orphan detection
- Stake-based validator system

---

## Use Cases

- **Audit Logs**: Immutable record of system events
- **Supply Chain**: Track product provenance
- **Document Verification**: Prove document existence and authenticity
- **Compliance**: Regulatory record-keeping
- **Voting Systems**: Transparent, verifiable elections
- **Academic Records**: Tamper-proof credential verification

---
