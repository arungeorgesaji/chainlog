package economy

import (
	"chainlog/core"
	"fmt"
)

type TransactionProcessor struct {
	FeeCollector string 
	BurnAddress  string 
}

func NewTransactionProcessor(feeCollector string) *TransactionProcessor {
	return &TransactionProcessor{
		FeeCollector: feeCollector,
		BurnAddress:  "BURN0000000000000000000000000000", 
	}
}

func (tp *TransactionProcessor) ProcessTransaction(tx *core.Transaction) error {
	if !ValidateTransactionFee(tx.Fee) {
		return fmt.Errorf("invalid transaction fee: %d (must be %d-%d)", 
			tx.Fee, MinTransactionFee, MaxTransactionFee)
	}
	
	if tx.Amount > 0 && tx.Receiver == "" {
		return fmt.Errorf("transfer transaction must have a receiver")
	}
	
	fmt.Printf("Transaction processed: %s paid %d LogCoin fee\n", 
		tx.Sender[:8], tx.Fee)
	
	return nil
}

func (tp *TransactionProcessor) CalculateFeeDistribution(fee uint64) (uint64, uint64) {
	minerShare := fee * 80 / 100
	burnAmount := fee * 20 / 100
	
	return minerShare, burnAmount
}

func (tp *TransactionProcessor) CreateFeeDistributionTransactions(tx *core.Transaction) []*core.Transaction {
	minerShare, burnAmount := tp.CalculateFeeDistribution(tx.Fee)
	
	var feeTxs []*core.Transaction
	
	if minerShare > 0 {
		feeTxs = append(feeTxs, &core.Transaction{
			ID:        fmt.Sprintf("fee-%s", tx.ID),
			Type:      core.FeeTx,
			Sender:    tx.Sender,
			Receiver:  tp.FeeCollector,
			Amount:    minerShare,
			Fee:       0, 
			Timestamp: tx.Timestamp,
			Nonce:     tx.Nonce + 1, 
		})
	}
	
	if burnAmount > 0 {
		feeTxs = append(feeTxs, &core.Transaction{
			ID:        fmt.Sprintf("burn-%s", tx.ID),
			Type:      core.FeeTx,
			Sender:    tx.Sender,
			Receiver:  tp.BurnAddress,
			Amount:    burnAmount,
			Fee:       0,
			Timestamp: tx.Timestamp,
			Nonce:     tx.Nonce + 2,
		})
	}
	
	return feeTxs
}
