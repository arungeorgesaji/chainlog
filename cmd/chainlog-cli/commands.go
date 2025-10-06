package main

import (
	"chainlog/core"
	"chainlog/crypto"
	"chainlog/network"
	"chainlog/consensus"
	"chainlog/economy"
	"chainlog/storage"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	bc     *core.Blockchain
	state  *storage.StateManager
	node   *network.Node
	ledger *storage.LedgerManager
)

func startNode() {
	port := "8080"
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	if err := storage.SaveNodeState(port); err != nil {
		fmt.Printf("Warning: Could not save node state: %v\n", err)
	}

	fmt.Printf("Starting ChainLog node on port %s...\n", port)
	
	bc = core.NewBlockchain()
	state = storage.NewStateManager()
	ledger = storage.NewLedgerManager(bc)

	if err := ledger.LoadBlockchain(); err != nil {
		fmt.Printf("Starting fresh blockchain: %v\n", err)
	}
	if err := state.LoadState(); err != nil {
		fmt.Printf("Starting fresh state: %v\n", err)
	}

	wallet, err := loadOrCreateWallet()
	if err != nil {
		storage.DeleteNodeState() 
		panic(err)
	}

	node = network.NewNode("localhost:"+port, wallet, bc, true)

	go node.CheckBroadcastFile()
	
	if err := node.Start(); err != nil {
		storage.DeleteNodeState()
		panic(err)
	}

	fmt.Printf("Node started successfully! Address: %s\n", wallet.GetAddress())
	fmt.Println("Node is running... (Ctrl+C to stop)")

	defer storage.DeleteNodeState()

	select {}
}

func handleWallet() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: chainlog-cli wallet [create|import|list|delete|info|default]")
		fmt.Println("\nCommands:")
		fmt.Println("  create [label]    - Create new wallet")
		fmt.Println("  import <key> [label] - Import wallet from private key")
		fmt.Println("  list              - List all saved wallets")
		fmt.Println("  delete <address>  - Delete wallet by address")
		fmt.Println("  info <address>    - Show wallet details")
		fmt.Println("  default           - Show default wallet")
		return
	}

	wm := storage.GetWalletManager()

	switch os.Args[2] {
	case "create":
		label := "My Wallet"
		if len(os.Args) >= 4 {
			label = os.Args[3]
		}

		wallet, err := crypto.NewWallet()
		if err != nil {
			fmt.Printf("Error creating wallet: %v\n", err)
			return
		}

		if err := wm.SaveWallet(wallet, label); err != nil {
			fmt.Printf("Error saving wallet: %v\n", err)
			return
		}

		fmt.Printf("Wallet created and saved successfully!\n\n")
		wallet.Display()
		fmt.Printf("\nLabel: %s\n", label)
		fmt.Printf("Storage: %s\n", storage.DataDir+"/"+storage.WalletsFile)
		fmt.Printf("Total wallets: %d\n", wm.WalletCount())

	case "import":
		if len(os.Args) < 4 {
			fmt.Println("Usage: chainlog-cli wallet import <private-key> [label]")
			return
		}

		label := "Imported Wallet"
		if len(os.Args) >= 5 {
			label = os.Args[4]
		}

		wallet, err := crypto.WalletFromPrivateKey(os.Args[3])
		if err != nil {
			fmt.Printf("Error importing wallet: %v\n", err)
			return
		}

		if err := wm.SaveWallet(wallet, label); err != nil {
			fmt.Printf("Error saving wallet: %v\n", err)
			return
		}

		fmt.Printf("Wallet imported and saved successfully!\n\n")
		wallet.Display()
		fmt.Printf("\nLabel: %s\n", label)

	case "list":
		wallets := wm.GetAllWallets()
		if len(wallets) == 0 {
			fmt.Println("No wallets found.")
			fmt.Println("   Create one with: chainlog-cli wallet create")
			return
		}

		fmt.Printf("Stored Wallets (%d):\n\n", len(wallets))
		for i, wallet := range wallets {
			created := time.Unix(wallet.CreatedAt, 0).Format("2006-01-02 15:04")
			fmt.Printf("%d. %s\n", i+1, wallet.Address)
			fmt.Printf("   Label: %s\n", wallet.Label)
			fmt.Printf("   Created: %s\n", created)
			if i < len(wallets)-1 {
				fmt.Println("   ──────────────────────────")
			}
		}

	case "delete":
		if len(os.Args) < 4 {
			fmt.Println("Usage: chainlog-cli wallet delete <address>")
			return
		}

		address := os.Args[3]
		if wm.DeleteWallet(address) {
			fmt.Printf("Wallet %s... deleted successfully\n", address[:8])
		} else {
			fmt.Printf("Wallet not found: %s\n", address)
		}

	case "info":
		if len(os.Args) < 4 {
			fmt.Println("Usage: chainlog-cli wallet info <address>")
			return
		}

		address := os.Args[3]
		stored, exists := wm.GetWallet(address)
		if !exists {
			fmt.Printf("Wallet not found: %s\n", address)
			return
		}

		wallet, err := crypto.WalletFromPrivateKey(stored.PrivateKey)
		if err != nil {
			fmt.Printf("Error loading wallet: %v\n", err)
			return
		}

		fmt.Printf("Wallet Details:\n")
		fmt.Printf("├─ Address: %s\n", wallet.Address)
		fmt.Printf("├─ Label: %s\n", stored.Label)
		fmt.Printf("├─ Created: %s\n", time.Unix(stored.CreatedAt, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("└─ Public Key: %s...\n", stored.PublicKey[:16])

	case "default":
		wallet, err := storage.GetDefaultWallet()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("   Create a wallet first with: chainlog-cli wallet create")
			return
		}

		stored, _ := wm.GetWallet(wallet.Address)
		fmt.Printf("Default Wallet:\n")
		fmt.Printf("├─ Address: %s\n", wallet.Address)
		fmt.Printf("├─ Label: %s\n", stored.Label)
		fmt.Printf("└─ Short: %s\n", wallet.GetAddressShort())

	default:
		fmt.Println("Usage: chainlog-cli wallet [create|import|list|delete|info|default]")
	}
}

func handleTransaction() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: chainlog-cli transaction [create|list|broadcast|status]")
		fmt.Println("\nCommands:")
		fmt.Println("  create <data> <fee> [wallet_address] - Create new transaction")
		fmt.Println("  list                                 - List pending transactions")
		fmt.Println("  broadcast <tx_id>                    - Broadcast transaction")
		fmt.Println("  status <tx_id>                       - Check transaction status")
		return
	}

	switch os.Args[2] {
	case "create":
		handleTransactionCreate()
	case "list":
		handleTransactionList()
	case "broadcast":
		handleTransactionBroadcast()
	case "status":
		handleTransactionStatus()
	default:
		fmt.Println("Usage: chainlog-cli transaction [create|list|broadcast|status]")
	}
}

func handleTransactionCreate() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: chainlog-cli transaction create <data> <fee> [wallet_address]")
		fmt.Println("\nExamples:")
		fmt.Println("  chainlog-cli transaction create \"Research data: temp=25C\" 10")
		fmt.Println("  chainlog-cli transaction create \"System log: user login\" 5 0x1234...")
		return
	}

	data := os.Args[3]
	fee, err := strconv.ParseUint(os.Args[4], 10, 64)
	if err != nil {
		fmt.Printf("Invalid fee: %v\n", err)
		return
	}

	var wallet *crypto.Wallet
	if len(os.Args) >= 6 {
		walletAddress := os.Args[5]
		wallet, err = storage.LoadWalletFromStorage(walletAddress)
		if err != nil {
			fmt.Printf("Error loading wallet %s: %v\n", walletAddress[:8], err)
			return
		}
	} else {
		wallet, err = storage.GetDefaultWallet()
		if err != nil {
			fmt.Printf("No wallet found: %v\n", err)
			fmt.Println("   Create a wallet first: chainlog-cli wallet create")
			fmt.Println("   Or specify wallet: chainlog-cli transaction create <data> <fee> <wallet_address>")
			return
		}
	}

	tx, err := core.NewDataTransaction(data, wallet, fee)
	if err != nil {
		fmt.Printf("Error creating transaction: %v\n", err)
		return
	}

	bc.AddTransaction(tx)

	fmt.Printf("Transaction created successfully!\n\n")
	fmt.Printf("Transaction Details:\n")
	fmt.Printf("├─ ID: %s\n", tx.ID)
	fmt.Printf("├─ From: %s\n", wallet.GetAddressShort())
	fmt.Printf("├─ Data: %s\n", tx.Data)
	fmt.Printf("├─ Fee: %d LogCoins\n", tx.Fee)
	fmt.Printf("├─ Timestamp: %d\n", tx.Timestamp)
	fmt.Printf("└─ Status: Pending\n")

	if ledger != nil {
		if err := ledger.SaveBlockchain(); err != nil {
			fmt.Printf("Warning: Could not save blockchain: %v\n", err)
		}
	}

	f, err := os.OpenFile("pending_broadcasts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err == nil {
    defer f.Close()
    f.WriteString(tx.ID + "\n")
    fmt.Printf("Transaction queued for broadcast\n")
}
}

func handleTransactionList() {
	pending := bc.GetPendingTransactions()
	
	if len(pending) == 0 {
		fmt.Println("No pending transactions")
		return
	}

	fmt.Printf("Pending Transactions (%d):\n\n", len(pending))
	for i, tx := range pending {
		fmt.Printf("%d. %s\n", i+1, tx.ID)
		fmt.Printf("   From: %s\n", tx.Sender[:8])
		fmt.Printf("   Data: %.50s\n", tx.Data)
		fmt.Printf("   Fee: %d LogCoins\n", tx.Fee)
		fmt.Printf("   Size: %d bytes\n", len(tx.Data))
		
		if i < len(pending)-1 {
			fmt.Println("   ──────────────────────────────────")
		}
	}
}

func handleTransactionStatus() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: chainlog-cli transaction status <tx_id>")
		return
	}

	txID := os.Args[3]
	
	pending := bc.GetPendingTransactions()
	for _, tx := range pending {
		if tx.ID == txID || strings.HasPrefix(tx.ID, txID) {
			fmt.Printf("Transaction Status:\n")
			fmt.Printf("├─ ID: %s...\n", tx.ID[:16])
			fmt.Printf("├─ Status: Pending\n")
			fmt.Printf("├─ Data: %.50s\n", tx.Data)
			fmt.Printf("├─ Fee: %d LogCoins\n", tx.Fee)
			fmt.Printf("└─ Created: %s\n", time.Unix(tx.Timestamp, 0).Format("2006-01-02 15:04:05"))
			return
		}
	}

	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.ID == txID || strings.HasPrefix(tx.ID, txID) {
				fmt.Printf("Transaction Status:\n")
				fmt.Printf("├─ ID: %s...\n", tx.ID[:16])
				fmt.Printf("├─ Status: Confirmed\n")
				fmt.Printf("├─ Block: %d\n", block.Index)
				fmt.Printf("├─ Data: %.50s\n", tx.Data)
				fmt.Printf("├─ Fee: %d LogCoins\n", tx.Fee)
				fmt.Printf("└─ Confirmed: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
				return
			}
		}
	}

	fmt.Printf("Transaction not found: %s\n", txID)
}

func handleTransactionBroadcast() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: chainlog-cli transaction broadcast <tx_id>")
		return
	}

	txID := os.Args[3]
	
	if !storage.IsNodeRunning() {
		fmt.Println("No node is currently running!")
		fmt.Println("   Start a node first: chainlog-cli start <port>")
		return
	}

	f, err := os.OpenFile("pending_broadcasts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creating broadcast file: %v\n", err)
		return
	}
	defer f.Close()
	
	if _, err := f.WriteString(txID + "\n"); err != nil {
		fmt.Printf("Error writing to broadcast file: %v\n", err)
		return
	}

	fmt.Printf("Transaction %s... queued for broadcast\n", txID[:16])
	fmt.Printf("   Node will pick it up within 10 seconds\n")
	fmt.Printf("   Check node terminal for broadcast confirmation\n")
}

func handleMine() {
    fmt.Println("Mining pending transactions...")
    
    wallet, err := loadOrCreateWallet()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    miner := consensus.NewMiner(wallet, bc)
    block, err := miner.MineBlock()
    if err != nil {
        fmt.Printf("Mining failed: %v\n", err)
        return
    }

    bc.Chain = append(bc.Chain, block)
    bc.ClearPendingTransactions()

    if ledger != nil {
        if err := ledger.SaveBlockchain(); err != nil {
            fmt.Printf("Warning: Could not save blockchain: %v\n", err)
        }
    }
    if state != nil {
        if err := state.SaveState(); err != nil {
            fmt.Printf("Warning: Could not save state: %v\n", err)
        }
    }

    fmt.Printf("Successfully mined block %d!\n", block.Index)
    fmt.Printf("Block hash: %s\n", block.Hash[:16])
}

func handleStatus() {
	fmt.Printf("\nCHAINLOG STATUS\n")
	fmt.Printf("├─ Blocks: %d\n", bc.GetBlockCount())
	fmt.Printf("├─ Pending Transactions: %d\n", len(bc.GetPendingTransactions()))
	fmt.Printf("├─ Difficulty: %d\n", bc.Difficulty)
	fmt.Printf("└─ Valid: %t\n", core.NewValidator(bc).ValidateBlockchain())

	if node != nil {
		fmt.Printf("├─ Node Peers: %d\n", node.GetPeerCount())
	}
}

func handleBalance() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: chainlog-cli balance <address>")
		return
	}

	address := os.Args[2]
	balance := state.GetBalance(address)
	fmt.Printf("Balance for %s: %d LogCoins\n", address[:8], balance)
}

func handlePeers() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: chainlog-cli peers [add|list]")
		return
	}

	switch os.Args[2] {
	case "add":
    if len(os.Args) < 4 {
        fmt.Println("Usage: chainlog-cli peers add <address>")
        return
    }
    fmt.Println("Peer-to-peer networking not fully implemented yet")

	case "list":
    if storage.IsNodeRunning() {
        fmt.Printf("Peers list unavailable across terminals as currently.\n")
        fmt.Printf("Check the node's terminal output or use the same terminal.\n")
    } else {
        fmt.Println("No node is running. Use: chainlog-cli start")
    }

	default:
		fmt.Println("Usage: chainlog-cli peers [add|list]")
	}
}

func loadOrCreateWallet() (*crypto.Wallet, error) {
	wm := storage.GetWalletManager()
	
	wallets := wm.GetAllWallets()
	
	if len(wallets) > 0 {
		fmt.Printf("Found %d existing wallet(s), using default...\n", len(wallets))
		wallet, err := storage.GetDefaultWallet()
		if err != nil {
			return nil, fmt.Errorf("failed to load default wallet: %v", err)
		}
		fmt.Printf("Loaded wallet: %s (%s)\n", wallet.GetAddressShort(), wallets[0].Label)
		return wallet, nil
	}
	
	fmt.Println("No existing wallets found, creating new wallet...")
	wallet, err := crypto.NewWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}
	
	label := "Node Wallet " + time.Now().Format("2006-01-02")
	if err := wm.SaveWallet(wallet, label); err != nil {
		return nil, fmt.Errorf("failed to save wallet: %v", err)
	}
	
	fmt.Printf("Created new wallet: %s (%s)\n", wallet.GetAddressShort(), label)
	return wallet, nil
}

func handleChain() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: chainlog-cli chain [show|validate]")
		return
	}

	switch os.Args[2] {
	case "show":
		handleChainShow()
	case "validate":
		handleChainValidate()
	default:
		fmt.Println("Usage: chainlog-cli chain [show|validate]")
	}
}

func handleChainShow() {
    fmt.Printf("\nBLOCKCHAIN (%d blocks)\n", len(bc.Chain))
    fmt.Println("═══════════════════════════════════════════════════")
    
    for i, block := range bc.Chain {
        fmt.Printf("Block %d:\n", block.Index)
        
        if len(block.Hash) >= 16 {
            fmt.Printf("├─ Hash: %s...\n", block.Hash[:16])
        } else if len(block.Hash) > 0 {
            fmt.Printf("├─ Hash: %s\n", block.Hash)
        } else {
            fmt.Printf("├─ Hash: (empty)\n")
        }
        
        if len(block.PrevHash) >= 16 {
            fmt.Printf("├─ Previous: %s...\n", block.PrevHash[:16])
        } else if len(block.PrevHash) > 0 {
            fmt.Printf("├─ Previous: %s\n", block.PrevHash)
        } else {
            fmt.Printf("├─ Previous: (genesis)\n")
        }
        
        fmt.Printf("├─ Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
        fmt.Printf("├─ Transactions: %d\n", len(block.Transactions))
        fmt.Printf("├─ Nonce: %d\n", block.Nonce)
        fmt.Printf("└─ Difficulty: %d\n", block.Difficulty)
        
        if i < len(bc.Chain)-1 {
            fmt.Println("   │")
        }
    }
}

func handleChainValidate() {
	validator := core.NewValidator(bc)
	isValid := validator.ValidateBlockchain()
	
	if isValid {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain validation failed!")
	}
}

func handleEconomyStats() {
	fmt.Println("\nCHAINLOG ECONOMICS")
	fmt.Println("═══════════════════════════════════════════════════")
	economy.DisplayEconomics()
}

func handleFeesStats() {
    totalFees := uint64(0)
    totalBurned := uint64(0)
    
    for _, block := range bc.Chain {
        for _, tx := range block.Transactions {
            if tx.Fee > 0 {
                totalFees += tx.Fee
                totalBurned += tx.Fee / 10
            }
        }
    }
    
    fmt.Printf("\nFEE STATISTICS (Real Blockchain Data)\n")
    fmt.Printf("├─ Total Fees Collected: %d LogCoins\n", totalFees)
    fmt.Printf("├─ Total Fees Burned: %d LogCoins\n", totalBurned)
    fmt.Printf("├─ Total Transactions: %d\n", countTransactionsWithFees())
    fmt.Printf("└─ Fee Efficiency: %.1f%%\n", float64(totalFees-totalBurned)/float64(totalFees)*100)
}

func countTransactionsWithFees() int {
    count := 0
    for _, block := range bc.Chain {
        for _, tx := range block.Transactions {
            if tx.Fee > 0 {
                count++
            }
        }
    }
    return count
}

func handleRewardsStats() {
    totalRewards := uint64(0)
    blocksMined := 0
    
    for _, block := range bc.Chain {
        for _, tx := range block.Transactions {
            if tx.Type == core.RewardTx {
                totalRewards += tx.Amount
                blocksMined++
            }
        }
    }
    
    wallet, _ := loadOrCreateWallet()
    
    fmt.Printf("\nREWARD STATISTICS (Real Blockchain Data)\n")
    fmt.Printf("├─ Total Rewards Distributed: %d LogCoins\n", totalRewards)
    fmt.Printf("├─ Blocks Mined: %d\n", blocksMined)
    fmt.Printf("├─ Average Reward per Block: %d LogCoins\n", totalRewards/uint64(max(blocksMined, 1)))
    fmt.Printf("└─ Your Miner: %s\n", wallet.GetAddressShort())
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func handleStaking() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: chainlog-cli staking [add|list]")
		return
	}

	switch os.Args[2] {
	case "add":
		handleStakingAdd()
	case "list":
		handleStakingList()
	default:
		fmt.Println("Usage: chainlog-cli staking [add|list]")
	}
}

func handleStakingAdd() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: chainlog-cli staking add <address> <amount>")
		return
	}

	address := os.Args[3]
	amount, err := strconv.ParseUint(os.Args[4], 10, 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		return
	}

	staking := consensus.NewStakingManager()
	staking.AddStake(address, amount)

	if err := staking.SaveStakingData(); err != nil {
		fmt.Printf("Warning: Could not save staking data: %v\n", err)
	}
	
	fmt.Printf("Added stake: %d LogCoins for %s...\n", amount, address[:8])
}

func handleStakingList() {
	staking := consensus.NewStakingManager()
	staking.DisplayValidators()
}

func handleDifficultyCheck() {
	diffManager := consensus.NewDifficultyManager(bc)
	newDifficulty := diffManager.CalculateNewDifficulty()
	
	fmt.Printf("\nDIFFICULTY ANALYSIS\n")
	fmt.Printf("├─ Current Difficulty: %d\n", bc.Difficulty)
	fmt.Printf("├─ Recommended Difficulty: %d\n", newDifficulty)
	fmt.Printf("└─ Adjustment: %+d\n", newDifficulty-bc.Difficulty)
}

func handleSave() {
	if ledger == nil {
		ledger = storage.NewLedgerManager(bc)
	}

	fmt.Println("Saving blockchain and state to disk...")
	
	if err := ledger.SaveBlockchain(); err != nil {
		fmt.Printf("Failed to save blockchain: %v\n", err)
	} else {
		fmt.Println(" Blockchain saved successfully!")
	}
	
	if err := state.SaveState(); err != nil {
		fmt.Printf("Failed to save state: %v\n", err)
	} else {
		fmt.Println("Account state saved successfully!")
	}
}

func handleLoad() {
	if ledger == nil {
		ledger = storage.NewLedgerManager(bc)
	}

	fmt.Println("Loading blockchain and state from disk...")
	
	if err := ledger.LoadBlockchain(); err != nil {
		fmt.Printf("Failed to load blockchain: %v\n", err)
	} else {
		fmt.Printf("Loaded blockchain with %d blocks\n", len(bc.Chain))
	}
	
	if err := state.LoadState(); err != nil {
		fmt.Printf("Failed to load state: %v\n", err)
	} else {
		fmt.Printf("Loaded state with %d accounts\n", len(state.Accounts))
	}
}

func handleSummary() {
	fmt.Println("\nCHAINLOG SYSTEM SUMMARY")
	fmt.Println("═══════════════════════════════════════════════════")
	
	fmt.Printf("Blockchain:\n")
	fmt.Printf("├─ Blocks: %d\n", bc.GetBlockCount())
	fmt.Printf("├─ Pending Transactions: %d\n", len(bc.GetPendingTransactions()))
	fmt.Printf("├─ Difficulty: %d\n", bc.Difficulty)
	fmt.Printf("└─ Valid: %t\n", core.NewValidator(bc).ValidateBlockchain())
	
	if node != nil {
		fmt.Printf("Network:\n")
		fmt.Printf("├─ Peers: %d\n", node.GetPeerCount())
		fmt.Printf("└─ Address: %s\n", node.Address)
	}
	
	wm := storage.GetWalletManager()
	fmt.Printf("Wallets:\n")
	fmt.Printf("├─ Total: %d\n", wm.WalletCount())
	
	fmt.Printf("Accounts:\n")
	fmt.Printf("└─ Tracked: %d\n", len(state.Accounts))
	
	fmt.Printf("Storage:\n")
	fmt.Printf("├─ Blockchain File: %s\n", storage.BlocksFile)
	fmt.Printf("├─ State File: %s\n", storage.StateFile)
	fmt.Printf("└─ Wallets File: %s\n", storage.WalletsFile)
	
	staking := consensus.NewStakingManager()
	fmt.Printf("Staking:\n")
	fmt.Printf("├─ Validators: %d\n", len(staking.Validators))
	fmt.Printf("└─ Total Staked: %d LogCoins\n", staking.GetTotalStaked())
}
