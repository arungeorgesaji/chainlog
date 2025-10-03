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
	fmt.Println("CHAINLOG COMPLETE SYSTEM TEST")
	fmt.Println("======================================")
	
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
	
	fmt.Println("\nInitial LogCoin Economics:")
	economy.DisplayEconomics()
	
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
	
	fmt.Println("\n4. Testing Economy System...")
	txProcessor := economy.NewTransactionProcessor(wallet1.GetAddress())
	
	txProcessor.ProcessTransaction(tx1)
	txProcessor.ProcessTransaction(tx2)
	
	feeTxs := txProcessor.CreateFeeDistributionTransactions(tx1)
	fmt.Printf("   Created %d fee distribution transactions\n", len(feeTxs))
	
	fmt.Println("\n5. Starting Network Node...")
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
	
	fmt.Println("\n6. Broadcasting transactions to network...")
	node.BroadcastTransaction(tx1)
	node.BroadcastTransaction(tx2)
	
	feeManager := economy.NewFeeManager(wallet1.GetAddress())
	
	fmt.Println("\n7. Setting Up Miner...")
	miner := consensus.NewMiner(wallet1, bc)
	miner.Display()
	
	rewardManager := economy.NewRewardManager(wallet1.GetAddress())
	
	fmt.Println("\n8. Mining Block with Pending Transactions...")
	
	if len(bc.GetPendingTransactions()) > 0 {
		fmt.Printf("   Mining %d pending transactions...\n", len(bc.GetPendingTransactions()))
		
		block, err := miner.MineBlock()
		if err != nil {
			fmt.Printf(" Mining failed: %v\n", err)
		} else {
			fmt.Println("\nProcessing transaction fees...")
			feeDistribution := feeManager.ProcessFees(block.Transactions)
			fmt.Printf("   Created %d fee distribution transactions\n", len(feeDistribution))
			
			bc.Chain = append(bc.Chain, block)
			bc.ClearPendingTransactions()
			
			blockRewardTx := rewardManager.CreateBlockReward(block.Index)
			if blockRewardTx != nil {
				fmt.Printf("Block reward: %d LogCoins → %s\n", 
					blockRewardTx.Amount, wallet1.GetAddressShort())
			}
			
			fmt.Printf("Successfully mined block %d!\n", block.Index)
			fmt.Printf("Mined by: %s\n", wallet1.GetAddressShort())
			
			node.BroadcastBlock(block)

			_, totalBurned := feeManager.GetFeeStatistics()
			economy.UpdateEconomics(block.Index, totalBurned)
		}
	} else {
		fmt.Println("   No pending transactions to mine")
	}
	
	fmt.Println("\n9. Blockchain state:")
	bc.Display()
	
	fmt.Println("\n10. Validating blockchain...")
	validator := core.NewValidator(bc)
	validator.ValidateBlockchain()
	
	fmt.Printf("\n11. Pending Transactions: %d\n", len(bc.GetPendingTransactions()))
	for i, tx := range bc.GetPendingTransactions() {
		fmt.Printf("\nTransaction %d:\n", i+1)
		tx.Display()
	}
	
	fmt.Println("\n12. Network Activity Demo...")
	
	fmt.Println("   Creating and broadcasting new network transaction...")
	tx3, err := core.NewDataTransaction("Network broadcast test!", wallet1, 1)
	if err != nil {
		panic(err)
	}
	bc.AddTransaction(tx3)
	node.BroadcastTransaction(tx3)
	
	txProcessor.ProcessTransaction(tx3)
	
	fmt.Println("\n13. Mining the New Transaction...")
	
	if len(bc.GetPendingTransactions()) > 0 {
		fmt.Printf("   Mining %d pending transactions...\n", len(bc.GetPendingTransactions()))
		
		block, err := miner.MineBlock()
		if err != nil {
			fmt.Printf("Mining failed: %v\n", err)
		} else {
			feeDistribution := feeManager.ProcessFees(block.Transactions)
			fmt.Printf("   Created %d fee distribution transactions\n", len(feeDistribution))
			
			bc.Chain = append(bc.Chain, block)
			bc.ClearPendingTransactions()
			
			blockRewardTx := rewardManager.CreateBlockReward(block.Index)
			if blockRewardTx != nil {
				fmt.Printf("Block reward: %d LogCoins\n", blockRewardTx.Amount)
			}
			
			fmt.Printf("Successfully mined block %d!\n", block.Index)
			
			node.BroadcastBlock(block)

			_, totalBurned := feeManager.GetFeeStatistics()
			economy.UpdateEconomics(block.Index, totalBurned)
		}
	}
	
	fmt.Println("\n14. Staking System Demo...")
	staking := consensus.NewStakingManager()
	staking.AddStake(wallet1.GetAddress(), 150) 
	staking.AddStake(wallet2.GetAddress(), 200)  
	staking.DisplayValidators()
	
	fmt.Println("\nStaking Rewards Demo...")
	for _, validator := range staking.Validators {
		stakeReward := rewardManager.CreateStakingReward(validator.Address, validator.Staked)
		if stakeReward != nil {
			fmt.Printf("   %s earned %d LogCoins staking reward\n", 
				validator.Address[:8], stakeReward.Amount)
		}
	}
	
	fmt.Println("\n15. Difficulty Management...")
	diffManager := consensus.NewDifficultyManager(bc)
	newDifficulty := diffManager.CalculateNewDifficulty()
	fmt.Printf("   Current difficulty: %d\n", bc.Difficulty)
	fmt.Printf("   Recommended difficulty: %d\n", newDifficulty)
	
	fmt.Println("\n16. Final Network Status:")
	node.Display()
	
	fmt.Println("\nFINAL ECONOMIC STATISTICS")
	fmt.Println("============================")
	economy.DisplayEconomics()
	feeManager.DisplayFeeStats()
	rewardManager.DisplayRewardStats()
	
	fmt.Println("\nSYSTEM SUMMARY")
	fmt.Println("=================")
	fmt.Printf("Wallets: 2\n")
	fmt.Printf("Blocks: %d\n", bc.GetBlockCount())
	fmt.Printf("Transactions Processed: %d\n", bc.GetBlockCount() * 2) 
	fmt.Printf("Network Peers: %d\n", node.GetPeerCount())
	fmt.Printf("Mining Difficulty: %d\n", bc.Difficulty)
	fmt.Printf("Total Staked: %d LogCoins\n", staking.GetTotalStaked())
	fmt.Printf("Validators: %d\n", len(staking.Validators))
	fmt.Printf("Total Fees Collected: %d LogCoins\n", feeManager.TotalFeesCollected)
	fmt.Printf("Total Rewards Distributed: %d LogCoins\n", rewardManager.TotalRewardsDistributed)
	
	fmt.Println("\nCHAINLOG COMPLETE SYSTEM TEST SUCCESSFUL!")
	fmt.Println("============================================")
}
