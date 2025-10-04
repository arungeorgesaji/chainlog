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
)

var (
	bc    *core.Blockchain
	state *storage.StateManager
	node  *network.Node
)

func startNode() {
	port := "8080"
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	fmt.Printf("Starting ChainLog node on port %s...\n", port)
	
	bc = core.NewBlockchain()
	state = storage.NewStateManager()
	ledger := storage.NewLedgerManager(bc)

	if err := ledger.LoadBlockchain(); err != nil {
		fmt.Printf("Starting fresh blockchain: %v\n", err)
	}
	if err := state.LoadState(); err != nil {
		fmt.Printf("Starting fresh state: %v\n", err)
	}

	wallet, err := loadOrCreateWallet()
	if err != nil {
		panic(err)
	}

	node = network.NewNode("localhost:"+port, wallet, bc, true)
	if err := node.Start(); err != nil {
		panic(err)
	}

	fmt.Printf("Node started successfully! Address: %s\n", wallet.GetAddress())
	fmt.Println("Node is running... (Ctrl+C to stop)")

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

	wm := crypto.GetWalletManager()

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
		wallet.DisplayFull()
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
		wallet, err := crypto.GetDefaultWallet()
		if err != nil {
			fmt.Printf("g %v\n", err)
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

func handleMine() {
	if bc == nil || state == nil {
		fmt.Println("Please start a node first")
		return
	}

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

	fmt.Printf("Successfully mined block %d!\n", block.Index)
	fmt.Printf("Block hash: %s\n", block.Hash[:16])
}

func handleStatus() {
	if bc == nil {
		fmt.Println("Please start a node first")
		return
	}

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
	if state == nil {
		fmt.Println("Please start a node first")
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
		if node == nil {
			fmt.Println("Please start a node first")
			return
		}
		node.AddPeer(os.Args[3])
		fmt.Printf("Added peer: %s\n", os.Args[3])

	case "list":
		if node == nil {
			fmt.Println("Please start a node first")
			return
		}
		node.Display()

	default:
		fmt.Println("Usage: chainlog-cli peers [add|list]")
	}
}

func loadOrCreateWallet() (*crypto.Wallet, error) {
    wm := crypto.GetWalletManager()
    
    wallets := wm.GetAllWallets()
    
    if len(wallets) > 0 {
        fmt.Printf("Found %d existing wallet(s), using default...\n", len(wallets))
        wallet, err := crypto.GetDefaultWallet()
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
		wallet, err = crypto.LoadWalletFromStorage(walletAddress)
		if err != nil {
			fmt.Printf("Error loading wallet %s: %v\n", walletAddress[:8], err)
			return
		}
	} else {
		wallet, err = crypto.GetDefaultWallet()
		if err != nil {
			fmt.Printf("No wallet found: %v\n", err)
			fmt.Println("   Create a wallet first: chainlog-cli wallet create")
			fmt.Println("   Or specify wallet: chainlog-cli transaction create <data> <fee> <wallet_address>")
			return
		}
	}

	if bc == nil {
		fmt.Println("Blockchain not initialized. Please start a node first.")
		fmt.Println("   Run: chainlog-cli start <port>")
		return
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

	if node != nil {
		node.BroadcastTransaction(tx)
		fmt.Printf("Transaction broadcast to network\n")
	}
}

func handleTransactionList() {
	if bc == nil {
		fmt.Println("Blockchain not initialized. Please start a node first.")
		fmt.Println("   Run: chainlog-cli start <port>")
		return
	}

	pending := bc.GetPendingTransactions()
	
	if len(pending) == 0 {
		fmt.Println("No pending transactions")
		return
	}

	fmt.Printf("Pending Transactions (%d):\n\n", len(pending))
	for i, tx := range pending {
		fmt.Printf("%d. %s...\n", i+1, tx.ID[:16])
		fmt.Printf("   From: %s...\n", tx.Sender[:8])
		fmt.Printf("   Data: %.50s\n", tx.Data)
		fmt.Printf("   Fee: %d LogCoins\n", tx.Fee)
		fmt.Printf("   Size: %d bytes\n", len(tx.Data))
		
		if i < len(pending)-1 {
			fmt.Println("   ──────────────────────────────────")
		}
	}
}

func handleTransactionBroadcast() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: chainlog-cli transaction broadcast <tx_id>")
		return
	}

	txID := os.Args[3]
	
	if bc == nil {
		fmt.Println("Blockchain not initialized.")
		return
	}

	if node == nil {
		fmt.Println("Node not running. Please start a node first.")
		return
	}

	pending := bc.GetPendingTransactions()
	var targetTx *core.Transaction
	
	for _, tx := range pending {
		if tx.ID == txID || strings.HasPrefix(tx.ID, txID) {
			targetTx = tx
			break
		}
	}

	if targetTx == nil {
		fmt.Printf("Transaction not found: %s\n", txID)
		return
	}

	node.BroadcastTransaction(targetTx)
	fmt.Printf("Transaction broadcast to network: %s...\n", targetTx.ID[:16])
}

func handleTransactionStatus() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: chainlog-cli transaction status <tx_id>")
		return
	}

	txID := os.Args[3]
	
	if bc == nil {
		fmt.Println("Blockchain not initialized.")
		return
	}

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

func handleMine() {
	if bc == nil || state == nil {
		fmt.Println("Please start a node first")
		return
	}

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

	fmt.Printf("Successfully mined block %d!\n", block.Index)
	fmt.Printf("Block hash: %s\n", block.Hash[:16])
}

func handleStatus() {
	if bc == nil {
		fmt.Println("Please start a node first")
		return
	}

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
	if state == nil {
		fmt.Println("Please start a node first")
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
		if node == nil {
			fmt.Println("Please start a node first")
			return
		}
		node.AddPeer(os.Args[3])
		fmt.Printf("Added peer: %s\n", os.Args[3])

	case "list":
		if node == nil {
			fmt.Println("Please start a node first")
			return
		}
		node.Display()

	default:
		fmt.Println("Usage: chainlog-cli peers [add|list]")
	}
}
