package main

import (
	"chainlog/core"
	"chainlog/crypto"
	"chainlog/network"
	"chainlog/consensus"
	"chainlog/economy"
	"chainlog/storage"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
		case "start":
			startNode()
		case "wallet":
			handleWallet()
		case "transaction":
			handleTransaction()
		case "mine":
			handleMine()
		case "status":
			handleStatus()
		case "balance":
			handleBalance()
		case "peers":
			handlePeers()
		case "help":
			printUsage()
		default:
			fmt.Println("Unknown command:", os.Args[1])
	}
}

func printUsage() {
	fmt.Println("ChainLog CLI")
	fmt.Println("==================================")
	fmt.Println("Commands:")
	fmt.Println("  start [port]                 - Start a node (default: 8080)")
	fmt.Println("  wallet create                 - Create a new wallet")
	fmt.Println("  wallet import <key>           - Import wallet from private key")
	fmt.Println("  wallet list                   - List all wallets")
	fmt.Println("  transaction create <data> <fee> - Create a transaction")
	fmt.Println("  transaction list              - List pending transactions")
	fmt.Println("  mine                          - Mine pending transactions")
	fmt.Println("  status                        - Show blockchain status")
	fmt.Println("  balance <address>             - Check account balance")
	fmt.Println("  peers add <address>           - Add a peer")
	fmt.Println("  peers list                    - List peers")
	fmt.Println("  chain show                    - Display full blockchain")
	fmt.Println("  chain validate                - Validate blockchain integrity")
	fmt.Println("  economy stats                 - Show LogCoin economics")
	fmt.Println("  fees stats                    - Show fee statistics")
	fmt.Println("  rewards stats                 - Show reward statistics")
	fmt.Println("  staking add <address> <amt>   - Stake LogCoins")
	fmt.Println("  staking list                  - List validators and stakes")
	fmt.Println("  difficulty check              - Show current vs. recommended difficulty")
	fmt.Println("  save                          - Save blockchain and state to disk")
	fmt.Println("  load                          - Load blockchain and state from disk")
	fmt.Println("  summary                       - Print full system summary")
}
