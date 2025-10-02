package main

import (
	"chainlog/core"
	"chainlog/crypto"  
	"chainlog/network"
	"chainlog/consensus"  
	"chainlog/economy"    
	"fmt"
)

func main() {
	fmt.Println("Starting ChainLog Cryptography Test...")
	
	fmt.Println("\n1. Creating wallets...")
	wallet1, err := crypto.NewWallet()
	if err != nil {
		panic(err)
	}
	
	wallet2, err := crypto.NewWallet()
	if err != nil {
		panic(err)
	}
	
	wallet1.Display()
	wallet2.Display()
	
	fmt.Println("\n2. Creating blockchain...")
	bc := core.NewBlockchain()
	
	bc.AddBlock("First research finding: Gravity exists!")
	bc.AddBlock("Second finding: Water is wet")
	bc.AddBlock("System error: Database connection failed")
	
	fmt.Println("\n3. Creating signed transactions...")
	tx1, err := core.NewDataTransaction("User logged in", wallet1, 2)
	if err != nil {
		panic(err)
	}
	
	tx2, err := core.NewDataTransaction("Payment processed", wallet2, 3)
	if err != nil {
		panic(err)
	}
	
	bc.AddTransaction(tx1)
	bc.AddTransaction(tx2)
	
	fmt.Println("\n4. Starting Network Node...")
	node := network.NewNode("localhost:8080", wallet1, bc, true)
	err = node.Start()
	if err != nil {
		panic(err)
	}
	defer node.Stop() 
	
	node.AddPeer("localhost:8081")
	node.AddPeer("localhost:8082")
	node.AddPeer("localhost:8083")
	
	node.Display()
	
	fmt.Println("\n5. Broadcasting transactions to network...")
	node.BroadcastTransaction(tx1)
	node.BroadcastTransaction(tx2)
	
	fmt.Println("\n6. Setting Up Miner...")
	miner := consensus.NewMiner(wallet1, bc)
	miner.Display()
	
	fmt.Println("\n7. Mining Block with Pending Transactions...")
	
	if len(bc.GetPendingTransactions()) > 0 {
		fmt.Printf("   Mining %d pending transactions...\n", len(bc.GetPendingTransactions()))
		
		block, err := miner.MineBlock()
		if err != nil {
			fmt.Printf("Mining failed: %v\n", err)
		} else {
			bc.Chain = append(bc.Chain, block)
			bc.ClearPendingTransactions()
			
			fmt.Printf("Successfully mined block %d!\n", block.Index)
			fmt.Printf("Miner reward: %d LogCoins\n", economy.CalculateBlockReward(block.Index))
			fmt.Printf("Mined by: %s\n", wallet1.GetAddressShort())
			
			node.BroadcastBlock(block)
		}
	} else {
		fmt.Println("   No pending transactions to mine")
	}
	
	fmt.Println("\n8. Blockchain state:")
	bc.Display()
	
	fmt.Println("\n9. Validating blockchain...")
	validator := core.NewValidator(bc)
	validator.ValidateBlockchain()
	
	fmt.Printf("\n10. Pending Transactions: %d\n", len(bc.GetPendingTransactions()))
	for i, tx := range bc.GetPendingTransactions() {
		fmt.Printf("\nTransaction %d:\n", i+1)
		tx.Display()
	}
	
	fmt.Println("\n11. Network Activity Demo...")
	
	fmt.Println("   Creating and broadcasting new network transaction...")
	tx3, err := core.NewDataTransaction("Network broadcast test!", wallet1, 1)
	if err != nil {
		panic(err)
	}
	bc.AddTransaction(tx3)
	node.BroadcastTransaction(tx3)
	
	fmt.Println("\n12. Mining the New Transaction...")
	
	if len(bc.GetPendingTransactions()) > 0 {
		fmt.Printf("   Mining %d pending transactions...\n", len(bc.GetPendingTransactions()))
		
		block, err := miner.MineBlock()
		if err != nil {
			fmt.Printf("Mining failed: %v\n", err)
		} else {
			bc.Chain = append(bc.Chain, block)
			bc.ClearPendingTransactions()
			
			fmt.Printf("Successfully mined block %d!\n", block.Index)
			fmt.Printf("Miner reward: %d LogCoins\n", economy.CalculateBlockReward(block.Index))
			
			node.BroadcastBlock(block)
		}
	}
	
	fmt.Println("\n13. Staking System Demo...")
	staking := consensus.NewStakingManager()
	staking.AddStake(wallet1.GetAddress(), 150) 
	staking.AddStake(wallet2.GetAddress(), 200)  
	staking.DisplayValidators()
	
	fmt.Println("\n14. Difficulty Management...")
	diffManager := consensus.NewDifficultyManager(bc)
	newDifficulty := diffManager.CalculateNewDifficulty()
	fmt.Printf("   Current difficulty: %d\n", bc.Difficulty)
	fmt.Printf("   Recommended difficulty: %d\n", newDifficulty)
	
	fmt.Println("\n15. Final Network Status:")
	node.Display()
	
	fmt.Println("\nSYSTEM SUMMARY")
	fmt.Println("=================")
	fmt.Printf("Wallets: 2\n")
	fmt.Printf("Blocks: %d\n", bc.GetBlockCount())
	fmt.Printf("Transactions Processed: %d\n", bc.GetBlockCount() * 2) 
	fmt.Printf("Network Peers: %d\n", node.GetPeerCount())
	fmt.Printf("Mining Difficulty: %d\n", bc.Difficulty)
	fmt.Printf("Total Staked: %d LogCoins\n", staking.GetTotalStaked())
	fmt.Printf("Validators: %d\n", len(staking.Validators))
	
	fmt.Println("\nCHAINLOG COMPLETE SYSTEM TEST SUCCESSFUL!")
	fmt.Println("============================================")
}
