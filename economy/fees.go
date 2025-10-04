package economy

import (
	"chainlog/core"
	"chainlog/storage"
	"fmt"
)

type FeeManager struct {
	TotalFeesCollected uint64
	TransactionProc    *TransactionProcessor 
	TotalFeesBurned    uint64
	FeeCollector       string
}

func NewFeeManager(feeCollector string, stateManager *storage.StateManager) *FeeManager {
	processor := NewTransactionProcessor(feeCollector, stateManager)

	return &FeeManager{
		FeeCollector:       feeCollector,
		TransactionProc:    processor,  
		TotalFeesCollected: 0,
		TotalFeesBurned:    0,
	}
}

func (fm *FeeManager) ProcessFees(transactions []*core.Transaction) []*core.Transaction {
	var feeDistributionTxs []*core.Transaction
	totalFees := uint64(0)
	totalBurned := uint64(0)
	
	processor := fm.TransactionProc
	
	for _, tx := range transactions {
		if tx.Fee > 0 {
			totalFees += tx.Fee
			
			minerShare, burnAmount := processor.CalculateFeeDistribution(tx.Fee)
			totalBurned += burnAmount
			
			feeTxs := processor.CreateFeeDistributionTransactions(tx)
			feeDistributionTxs = append(feeDistributionTxs, feeTxs...)
			
			fmt.Printf("   Fee: %d LogCoins → Miner: %d, Burned: %d\n", 
				tx.Fee, minerShare, burnAmount)
		}
	}
	
	fm.TotalFeesCollected += totalFees
	fm.TotalFeesBurned += totalBurned
	
	if totalFees > 0 {
		fmt.Printf("Block fees: %d collected, %d burned\n", totalFees, totalBurned)
	}
	
	return feeDistributionTxs
}

func (fm *FeeManager) GetFeeStatistics() (uint64, uint64) {
	return fm.TotalFeesCollected, fm.TotalFeesBurned
}

func (fm *FeeManager) DisplayFeeStats() {
	fmt.Printf("\nFEE STATISTICS\n")
	fmt.Printf("├─ Total Fees Collected: %d LogCoins\n", fm.TotalFeesCollected)
	fmt.Printf("├─ Total Fees Burned: %d LogCoins\n", fm.TotalFeesBurned)
	fmt.Printf("└─ Fee Collector: %s\n", fm.FeeCollector[:8])
}
