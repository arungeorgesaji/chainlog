package main

import (
	"chainlog/core"
	"chainlog/storage"
	"fmt"
	"os"
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
	case "help":
		printUsage()
	default:
		if !storage.IsNodeRunning() {
			fmt.Println("No node is currently running!")
			fmt.Println("   Start a node first: chainlog-cli start [port]")
			return
		}
		
		bc = core.NewBlockchain()
		state = storage.NewStateManager()
		ledger = storage.NewLedgerManager(bc)
		ledger.LoadBlockchain() 
		state.LoadState()       

		switch os.Args[1] {
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
		case "chain":
			handleChain()
		case "fees":
			handleFeesStats()
		case "rewards":
			handleRewardsStats()
		case "staking":
			handleStaking()
		case "difficulty":
			handleDifficultyCheck()
		case "save":
			handleSave()
		case "load":
			handleLoad()
		case "summary":
			handleSummary()
		default:
			fmt.Println("Unknown command:", os.Args[1])
			printUsage()
		}
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
	fmt.Println("  transaction broadcast <tx_id> - Broadcast transaction")
	fmt.Println("  transaction status <tx_id>    - Check transaction status")
	fmt.Println("  mine                          - Mine pending transactions")
	fmt.Println("  status                        - Show blockchain status")
	fmt.Println("  balance <address>             - Check account balance")
	fmt.Println("  peers add <address>           - Add a peer")
	fmt.Println("  peers list                    - List peers")
	fmt.Println("  chain show                    - Display full blockchain")
	fmt.Println("  chain validate                - Validate blockchain integrity")
	fmt.Println("  economy stats                 - Show LogCoin economics")
	fmt.Println("  fees                    - Show fee statistics")
	fmt.Println("  rewards                 - Show reward statistics")
	fmt.Println("  staking add <address> <amt>   - Stake LogCoins")
	fmt.Println("  staking list                  - List validators and stakes")
	fmt.Println("  difficulty check              - Show current vs. recommended difficulty")
	fmt.Println("  save                          - Save blockchain and state to disk")
	fmt.Println("  load                          - Load blockchain and state from disk")
	fmt.Println("  summary                       - Print full system summary")
	fmt.Println("  help                          - Show this help message")
}
