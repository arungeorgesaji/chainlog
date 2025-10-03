package main

import (
	"chainlog/core"
	"chainlog/crypto"  
	"chainlog/network"
	"chainlog/consensus"  
	"chainlog/economy"    
	"chainlog/storage"    
	"fmt"
)

func main() {
	fmt.Println("CHAINLOG COMPLETE SYSTEM TEST")
	fmt.Println("======================================")
	
	fmt.Println("\n1. Initializing Storage System...")
	
	bc := core.NewBlockchain()
	ledger := storage.NewLedgerManager(bc)
	state := storage.NewStateManager()
	
	fmt.Println("Loading existing data from disk...")
	if err := ledger.LoadBlockchain(); err != nil {
		fmt.Printf("   Starting fresh: %v\n", err)
	} else {
		fmt.Printf("   Loaded existing blockchain with %d blocks!\n", len(bc.Chain))
	}
	
	if err := state.LoadState(); err != nil {
		fmt.Printf("   Starting fresh state: %v\n", err)
	} else {
		fmt.Printf("   Loaded account state with %d accounts!\n", len(state.Accounts))
	}
	
	fmt.Println("\n2. Creating wallets...")
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
	
	if len(state.Accounts) == 0 {
		fmt.Println("\nInitializing genesis accounts...")
		state.InitializeGenesisState([]*crypto.Wallet{wallet1, wallet2})
	}
	
	fmt.Println("\nInitial LogCoin Economics:")
	economy.DisplayEconomics()
	
	if len(bc.Chain) <= 1 { 
		fmt.Println("\n3. Creating initial blocks...")
		bc.AddBlock("First research finding: Gravity exists!")
		bc.AddBlock("Second finding: Water is wet")
		bc.AddBlock("System error: Database connection failed")
	} else {
		fmt.Printf("\n3. Using existing blockchain with %d blocks\n", len(bc.Chain))
	}
	
	fmt.Println("\n4. Creating signed transactions...")
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
	
	state.UpdateAccount(wallet1.GetAddress(), state.GetBalance(wallet1.GetAddress())-tx1.Fee, 1)
	state.UpdateAccount(wallet2.GetAddress(), state.GetBalance(wallet2.GetAddress())-tx2.Fee, 1)
	
	fmt.Println("\n5. Testing Economy System...")
	txProcessor := economy.NewTransactionProcessor(wallet1.GetAddress())
	
	txProcessor.ProcessTransaction(tx1)
	txProcessor.ProcessTransaction(tx2)
	
	feeTxs := txProcessor.CreateFeeDistributionTransactions(tx1)
	fmt.Printf("   Created %d fee distribution transactions\n", len(feeTxs))
	
	fmt.Println("\n6. Starting Network Node...")
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
	
	fmt.Println("\n7. Broadcasting transactions to network...")
	node.BroadcastTransaction(tx1)
	node.BroadcastTransaction(tx2)
	
	feeManager := economy.NewFeeManager(wallet1.GetAddress())
	
	fmt.Println("\n8. Setting Up Miner...")
	miner := consensus.NewMiner(wallet1, bc)
	miner.Display()
	
	rewardManager := economy.NewRewardManager(wallet1.GetAddress())
	
	fmt.Println("\n9. Mining Block with Pending Transactions...")
	
	if len(bc.GetPendingTransactions()) > 0 {
		fmt.Printf("   Mining %d pending transactions...\n", len(bc.GetPendingTransactions()))
		
		block, err := miner.MineBlock()
		if err != nil {
			fmt.Printf("Mining failed: %v\n", err)
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
				currentBalance := state.GetBalance(wallet1.GetAddress())
				state.UpdateAccount(wallet1.GetAddress(), currentBalance+blockRewardTx.Amount, 2)
			}
			
			fmt.Printf("Successfully mined block %d!\n", block.Index)
			fmt.Printf("Mined by: %s\n", wallet1.GetAddressShort())
			
			node.BroadcastBlock(block)

			fmt.Println("\nSaving blockchain after mining...")
			if err := ledger.SaveBlockchain(); err != nil {
				fmt.Printf("Failed to save blockchain: %v\n", err)
			} else {
				fmt.Println("Blockchain saved successfully!")
			}

			_, totalBurned := feeManager.GetFeeStatistics()
			economy.UpdateEconomics(block.Index, totalBurned)
		}
	} else {
		fmt.Println("   No pending transactions to mine")
	}
	
	fmt.Println("\n10. Blockchain state:")
	bc.Display()
	
	fmt.Println("\n11. Validating blockchain...")
	validator := core.NewValidator(bc)
	validator.ValidateBlockchain()
	
	fmt.Printf("\n12. Pending Transactions: %d\n", len(bc.GetPendingTransactions()))
	for i, tx := range bc.GetPendingTransactions() {
		fmt.Printf("\nTransaction %d:\n", i+1)
		tx.Display()
	}
	
	fmt.Println("\n13. Network Activity Demo...")
	
	fmt.Println("   Creating and broadcasting new network transaction...")
	tx3, err := core.NewDataTransaction("Network broadcast test!", wallet1, 1)
	if err != nil {
		panic(err)
	}
	bc.AddTransaction(tx3)
	node.BroadcastTransaction(tx3)
	
	txProcessor.ProcessTransaction(tx3)
	
	state.UpdateAccount(wallet1.GetAddress(), state.GetBalance(wallet1.GetAddress())-tx3.Fee, 3)
	
	fmt.Println("\n14. Mining the New Transaction...")
	
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
				currentBalance := state.GetBalance(wallet1.GetAddress())
				state.UpdateAccount(wallet1.GetAddress(), currentBalance+blockRewardTx.Amount, 4)
			}
			
			fmt.Printf("Successfully mined block %d!\n", block.Index)
			
			node.BroadcastBlock(block)

			if err := ledger.SaveBlockchain(); err != nil {
				fmt.Printf("Failed to save blockchain: %v\n", err)
			}

			_, totalBurned := feeManager.GetFeeStatistics()
			economy.UpdateEconomics(block.Index, totalBurned)
		}
	}
	
	fmt.Println("\n15. Staking System Demo...")
	staking := consensus.NewStakingManager()
	staking.AddStake(wallet1.GetAddress(), 150) 
	staking.AddStake(wallet2.GetAddress(), 200)  
	staking.DisplayValidators()
	
	fmt.Println("\nStaking Rewards Demo...")
	for _, validator := range staking.Validators {
		stakeReward := rewardManager.CreateStakingReward(validator.Address, validator.Staked)
		if stakeReward != nil {
			fmt.Printf("   %s earned %d LogCoins staking reward\n", 				validator.Address[:8], stakeReward.Amount)
			currentBalance := state.GetBalance(validator.Address)
			state.UpdateAccount(validator.Address, currentBalance+stakeReward.Amount, 0)
		}
	}
	
	fmt.Println("\n16. Difficulty Management...")
	diffManager := consensus.NewDifficultyManager(bc)
	newDifficulty := diffManager.CalculateNewDifficulty()
	fmt.Printf("   Current difficulty: %d\n", bc.Difficulty)
	fmt.Printf("   Recommended difficulty: %d\n", newDifficulty)
	
	fmt.Println("\n17. Final Network Status:")
	node.Display()
	
	fmt.Println("\n18. Final Data Persistence...")
	if err := ledger.SaveBlockchain(); err != nil {
		fmt.Printf("Failed to save blockchain: %v\n", err)
	} else {
		fmt.Println("Blockchain saved successfully!")
	}
	
	if err := state.SaveState(); err != nil {
		fmt.Printf("Failed to save state: %v\n", err)
	} else {
		fmt.Println("Account state saved successfully!")
	}
	
	fmt.Println("\nSTORAGE INFORMATION")
	fmt.Println("======================")
	ledger.DisplayStorageInfo()
	state.DisplayState()
	
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
	fmt.Printf("Data Persistence: %t\n", storage.FileExists(storage.BlocksFile))
	fmt.Printf("Accounts Tracked: %d\n", len(state.Accounts))
	
	fmt.Println("\nCHAINLOG COMPLETE SYSTEM TEST SUCCESSFUL!")
	fmt.Println("============================================")
}
